#!/usr/bin/env bash
set -euo pipefail

# =============================================================================
# scripts/discuss.sh - 複数AI議論スクリプト
#
# Claude Haiku がファシリテータとなり、Claude Opus と Codex が
# チャット形式で技術的なトピックについて議論する。
# =============================================================================

readonly SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
readonly PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
readonly DISCUSSIONS_DIR="$PROJECT_ROOT/docs/discussions"
readonly SPECS_DIR="$PROJECT_ROOT/docs/specification"
readonly DEFAULT_MAX_ROUNDS=5
readonly MAX_CONTEXT_CHARS=30000
readonly CONVERGENCE_KEYWORD="CONSENSUS_REACHED"

# 色定義
readonly COLOR_RESET='\033[0m'
readonly COLOR_HAIKU='\033[0;36m'     # Cyan
readonly COLOR_OPUS='\033[0;35m'      # Magenta
readonly COLOR_CODEX='\033[0;33m'     # Yellow
readonly COLOR_HEADER='\033[1;37m'    # Bold White
readonly COLOR_DIM='\033[0;90m'       # Gray
readonly COLOR_ERROR='\033[0;31m'     # Red

# グローバル変数
TOPIC=""
CONTEXT_FILES=()
MAX_ROUNDS=$DEFAULT_MAX_ROUNDS
OUTPUT_FILE=""
CONVERSATION_HISTORY=""
CONTEXT=""

# システムプロンプト（init_session後に構築）
HAIKU_SYSTEM_PROMPT=""
OPUS_SYSTEM_PROMPT=""
CODEX_ROLE_PROMPT=""

# =============================================================================
# ユーティリティ
# =============================================================================

emit() {
    local speaker="$1"
    local message="$2"
    local color timestamp

    timestamp=$(date +"%H:%M:%S")

    case "$speaker" in
        HAIKU*)       color="$COLOR_HAIKU" ;;
        OPUS*)        color="$COLOR_OPUS" ;;
        CODEX*)       color="$COLOR_CODEX" ;;
        SYSTEM|WARNING) color="$COLOR_DIM" ;;
        ERROR)        color="$COLOR_ERROR" ;;
        *)            color="$COLOR_RESET" ;;
    esac

    # ターミナル出力（色付き）
    echo -e "${COLOR_DIM}[${timestamp}]${COLOR_RESET} ${color}[${speaker}]${COLOR_RESET}"
    echo -e "$message"
    echo ""

    # ファイル出力（Markdown）
    if [ -n "$OUTPUT_FILE" ]; then
        {
            echo "#### ${speaker} (${timestamp})"
            echo ""
            echo "$message"
            echo ""
            echo "---"
            echo ""
        } >> "$OUTPUT_FILE"
    fi
}

emit_header() {
    local title="$1"
    local separator="━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    echo -e "\n${COLOR_HEADER}${separator}${COLOR_RESET}"
    echo -e "${COLOR_HEADER}  ${title}${COLOR_RESET}"
    echo -e "${COLOR_HEADER}${separator}${COLOR_RESET}\n"

    if [ -n "$OUTPUT_FILE" ]; then
        echo -e "\n## ${title}\n" >> "$OUTPUT_FILE"
    fi
}

append_history() {
    local speaker="$1"
    local message="$2"

    CONVERSATION_HISTORY+="
### ${speaker}
${message}

"
}

# =============================================================================
# 引数パース
# =============================================================================

parse_args() {
    if [ $# -eq 0 ]; then
        echo "使い方: $0 <トピック> [--context <file>...] [--max-rounds <N>]"
        echo ""
        echo "例:"
        echo "  $0 \"バトルシステムのステートマシン設計\""
        echo "  $0 \"リファクタリング方針\" --context internal/battle/model.go"
        echo "  $0 \"APIの設計方針\" --max-rounds 3"
        exit 1
    fi

    TOPIC="$1"
    shift

    while [ $# -gt 0 ]; do
        case "$1" in
            --context)
                shift
                while [ $# -gt 0 ] && [[ "$1" != --* ]]; do
                    CONTEXT_FILES+=("$1")
                    shift
                done
                ;;
            --max-rounds)
                shift
                if [ $# -eq 0 ]; then
                    echo "エラー: --max-rounds に値が必要です"
                    exit 1
                fi
                MAX_ROUNDS="$1"
                shift
                ;;
            *)
                echo "不明なオプション: $1"
                exit 1
                ;;
        esac
    done
}

# =============================================================================
# セッション初期化
# =============================================================================

init_session() {
    local session_id
    session_id=$(date +"%Y%m%d-%H%M%S")

    # スラッグ生成（英数字とひらがな・カタカナ・漢字を保持）
    local slug
    slug=$(echo "$TOPIC" | tr ' ' '-' | head -c 50)

    mkdir -p "$DISCUSSIONS_DIR"
    OUTPUT_FILE="${DISCUSSIONS_DIR}/${session_id}-discussion.md"

    # ファイルヘッダ書き込み
    {
        echo "# 議論: ${TOPIC}"
        echo ""
        echo "**日時**: $(date '+%Y-%m-%d %H:%M:%S')"
        echo "**参加者**: Claude Opus (参加者), Codex (参加者), Claude Haiku (ファシリテータ)"
        echo "**最大ラウンド数**: ${MAX_ROUNDS}"
        echo ""
        if [ "${#CONTEXT_FILES[@]}" -gt 0 ]; then
            echo "## コンテキスト"
            for f in "${CONTEXT_FILES[@]+"${CONTEXT_FILES[@]}"}"; do
                echo "- \`$f\`"
            done
            echo ""
        fi
        echo "---"
        echo ""
    } > "$OUTPUT_FILE"

    emit "SYSTEM" "議論を開始します: $TOPIC"
    emit "SYSTEM" "出力先: $OUTPUT_FILE"
}

# =============================================================================
# コンテキスト収集
# =============================================================================

collect_context() {
    CONTEXT=""

    # 指定ファイルの読み込み
    for file in "${CONTEXT_FILES[@]+"${CONTEXT_FILES[@]}"}"; do
        local filepath
        if [[ "$file" = /* ]]; then
            filepath="$file"
        else
            filepath="$PROJECT_ROOT/$file"
        fi

        if [ -f "$filepath" ]; then
            CONTEXT+="
--- ${file} ---
$(cat "$filepath")

"
        else
            emit "WARNING" "ファイルが見つかりません: $file"
        fi
    done

    # トピックのキーワードから関連仕様を自動追加
    local spec_files=("tech.md" "structure.md")

    case "$TOPIC" in
        *バトル*|*battle*|*combat*)       spec_files+=("battle.md") ;;
        *エージェント*|*agent*|*スロット*) spec_files+=("agent.md") ;;
        *タイピング*|*typing*)            spec_files+=("typing.md") ;;
        *敵*|*enemy*)                     spec_files+=("enemy.md") ;;
        *図鑑*|*実績*|*collection*)       spec_files+=("collection.md") ;;
        *ゲームループ*|*gameloop*|*シーン*) spec_files+=("gameloop.md") ;;
    esac

    for spec in "${spec_files[@]}"; do
        local spec_path="$SPECS_DIR/$spec"
        if [ -f "$spec_path" ]; then
            CONTEXT+="
--- docs/specification/${spec} ---
$(cat "$spec_path")

"
        fi
    done

    # 長すぎる場合は切り詰め
    if [ ${#CONTEXT} -gt $MAX_CONTEXT_CHARS ]; then
        emit "WARNING" "コンテキストが上限 (${MAX_CONTEXT_CHARS}文字) を超えています。切り詰めます。"
        CONTEXT="${CONTEXT:0:$MAX_CONTEXT_CHARS}

... (以下、文字数上限により省略)"
    fi
}

# =============================================================================
# システムプロンプト構築
# =============================================================================

build_system_prompts() {
    HAIKU_SYSTEM_PROMPT="あなたは技術議論のファシリテータです。
2名の参加者（Claude Opus と Codex）の議論を整理・収束させる役割です。

## 責務
1. 各参加者の意見の共通点と相違点を明確にする
2. 議論が散漫にならないよう論点を絞る
3. 具体的な質問や検討事項を提示して議論を深める
4. 合意が形成されたと判断したら、その旨を明示する

## 出力構造

### 論点整理
{各参加者の主張のまとめと共通点/相違点}

### 深掘りポイント
{次に議論すべき具体的な質問。2-3個に絞る}

### 収束判定
{合意に至った場合のみ、行頭に CONSENSUS_REACHED と記載し、合意内容を箇条書きで記述する}

## 注意
- 自分の技術的意見は述べない（ファシリテーションに徹する）
- 参加者の意見を公平に扱う
- 議論が堂々巡りになったら、判断基準の明確化を促す
- 日本語で出力する"

    OPUS_SYSTEM_PROMPT="あなたはBlitzTypingOperatorプロジェクトの技術議論に参加するシニアエンジニアです。

## あなたの特性
- アーキテクチャ設計、設計パターン、長期的な保守性に重点を置く
- トレードオフを明示的に分析する
- 抽象的な議論よりも、プロジェクトの具体的なコードや構造に基づいて論じる

## 出力ルール
- 日本語で出力する
- 意見を述べる際は必ず根拠を示す
- 他の参加者の意見に対しては、同意/反対/補足を明確にする
- 300-500文字程度を目安に簡潔にまとめる"

    CODEX_ROLE_PROMPT="あなたはBlitzTypingOperatorプロジェクトの技術議論に参加するシニアエンジニアです。

## あなたの特性
- 実装の実現可能性、具体的なコードへの影響に重点を置く
- 既存コードベースとの整合性を重視する
- パフォーマンスやエッジケースに敏感

## 出力ルール
- 日本語で出力する
- 意見を述べる際は必ず根拠を示す
- 他の参加者の意見に対しては、同意/反対/補足を明確にする
- 300-500文字程度を目安に簡潔にまとめる"
}

# =============================================================================
# AI呼び出し
# =============================================================================

call_with_retry() {
    local max_retries=3
    local retry_delay=5
    local attempt=1

    while [ "$attempt" -le "$max_retries" ]; do
        local output
        if output=$("$@" 2>/dev/null); then
            echo "$output"
            return 0
        fi

        emit "ERROR" "呼び出し失敗 (試行 ${attempt}/${max_retries})"

        if [ "$attempt" -lt "$max_retries" ]; then
            emit "SYSTEM" "${retry_delay}秒後にリトライします..."
            sleep "$retry_delay"
            retry_delay=$((retry_delay * 2))
        fi

        attempt=$((attempt + 1))
    done

    emit "ERROR" "最大リトライ回数に達しました"
    return 1
}

call_haiku() {
    local prompt="$1"

    echo "$prompt" | call_with_retry claude -p \
        --model haiku \
        --system-prompt "$HAIKU_SYSTEM_PROMPT" \
        --no-session-persistence
}

call_opus() {
    local prompt="$1"

    echo "$prompt" | call_with_retry claude -p \
        --model opus \
        --system-prompt "$OPUS_SYSTEM_PROMPT" \
        --no-session-persistence
}

call_codex() {
    local prompt="$1"

    # codexはシステムプロンプトオプションがないため、役割定義をプロンプトに埋め込む
    local full_prompt="${CODEX_ROLE_PROMPT}

${prompt}"

    echo "$full_prompt" | call_with_retry codex exec --full-auto -
}

# =============================================================================
# プロンプト構築
# =============================================================================

build_initial_prompt() {
    cat <<EOF
以下の技術的トピックについて、あなたの意見を述べてください。

## トピック
${TOPIC}

## プロジェクトコンテキスト
${CONTEXT}

## 指示
- このプロジェクトの文脈で、トピックに対するあなたの立場と根拠を述べてください
- 複数のアプローチがある場合は、推奨するものを1つ選び、その理由を説明してください
- 300-500文字程度で簡潔にまとめてください
EOF
}

build_facilitation_prompt() {
    local round="$1"

    cat <<EOF
以下の技術的議論のファシリテーションを行ってください。

## トピック
${TOPIC}

## 議論経過
${CONVERSATION_HISTORY}

## 現在のラウンド
${round} / ${MAX_ROUNDS}

## 指示
- 各参加者の意見を整理し、共通点と相違点を明確にしてください
- 議論を深めるための具体的な質問を2-3個提示してください
- 参加者の意見が実質的に合意に達していると判断した場合は、収束判定セクションに CONSENSUS_REACHED と記載してください
EOF
}

build_response_prompt() {
    local facilitation="$1"

    cat <<EOF
以下の議論に対して応答してください。

## トピック
${TOPIC}

## これまでの議論
${CONVERSATION_HISTORY}

## ファシリテータからの問いかけ
${facilitation}

## 指示
- ファシリテータの質問に具体的に回答してください
- 他の参加者の意見への同意/反対/補足を明確にしてください
- 300-500文字程度で簡潔にまとめてください
EOF
}

build_conclusion_prompt() {
    cat <<EOF
以下の議論をまとめてください。

## トピック
${TOPIC}

## 議論の全履歴
${CONVERSATION_HISTORY}

## 指示
以下の構造でまとめを作成してください:

### 結論
{合意に至った決定事項。箇条書き}

### 根拠
{決定に至った主な理由}

### 残課題
{合意に至らなかった点、今後検討が必要な点}

### アクションアイテム
{具体的な次のステップ}
EOF
}

# =============================================================================
# 議論フェーズ
# =============================================================================

phase_initial_opinions() {
    emit_header "Phase 1: 初期意見"

    local prompt
    prompt=$(build_initial_prompt)

    local opus_tmpfile codex_tmpfile
    opus_tmpfile=$(mktemp)
    codex_tmpfile=$(mktemp)

    # 並列呼び出し
    call_opus "$prompt" > "$opus_tmpfile" &
    local opus_pid=$!

    call_codex "$prompt" > "$codex_tmpfile" &
    local codex_pid=$!

    # Opus の結果を待機
    local opus_output=""
    if wait "$opus_pid" 2>/dev/null; then
        opus_output=$(cat "$opus_tmpfile")
    else
        opus_output="(応答を取得できませんでした)"
        emit "ERROR" "Opus の呼び出しに失敗しました"
    fi

    # Codex の結果を待機
    local codex_output=""
    if wait "$codex_pid" 2>/dev/null; then
        codex_output=$(cat "$codex_tmpfile")
    else
        codex_output="(応答を取得できませんでした)"
        emit "ERROR" "Codex の呼び出しに失敗しました"
    fi

    rm -f "$opus_tmpfile" "$codex_tmpfile"

    emit "OPUS" "$opus_output"
    emit "CODEX" "$codex_output"

    append_history "OPUS" "$opus_output"
    append_history "CODEX" "$codex_output"
}

phase_facilitate() {
    local round="$1"

    local prompt
    prompt=$(build_facilitation_prompt "$round")

    local facilitation
    if facilitation=$(call_haiku "$prompt"); then
        emit "HAIKU (ファシリテータ)" "$facilitation"
        append_history "FACILITATOR" "$facilitation"
        echo "$facilitation"
    else
        emit "ERROR" "ファシリテータの呼び出しに失敗しました。議論を中断します。"
        return 1
    fi
}

phase_respond() {
    local facilitation="$1"

    local prompt
    prompt=$(build_response_prompt "$facilitation")

    local opus_tmpfile codex_tmpfile
    opus_tmpfile=$(mktemp)
    codex_tmpfile=$(mktemp)

    # 並列呼び出し
    call_opus "$prompt" > "$opus_tmpfile" &
    local opus_pid=$!

    call_codex "$prompt" > "$codex_tmpfile" &
    local codex_pid=$!

    # Opus の結果
    local opus_output=""
    if wait "$opus_pid" 2>/dev/null; then
        opus_output=$(cat "$opus_tmpfile")
    else
        opus_output="(応答を取得できませんでした)"
        emit "ERROR" "Opus の呼び出しに失敗しました"
    fi

    # Codex の結果
    local codex_output=""
    if wait "$codex_pid" 2>/dev/null; then
        codex_output=$(cat "$codex_tmpfile")
    else
        codex_output="(応答を取得できませんでした)"
        emit "ERROR" "Codex の呼び出しに失敗しました"
    fi

    rm -f "$opus_tmpfile" "$codex_tmpfile"

    emit "OPUS" "$opus_output"
    emit "CODEX" "$codex_output"

    append_history "OPUS" "$opus_output"
    append_history "CODEX" "$codex_output"
}

check_convergence() {
    local facilitation="$1"
    echo "$facilitation" | grep -q "$CONVERGENCE_KEYWORD"
}

phase_conclude() {
    emit_header "結論"

    local prompt
    prompt=$(build_conclusion_prompt)

    local conclusion
    if conclusion=$(call_haiku "$prompt"); then
        emit "HAIKU (結論)" "$conclusion"
    else
        emit "ERROR" "結論まとめの生成に失敗しました"
    fi
}

summarize_history() {
    emit "SYSTEM" "コンテキスト圧縮中..."

    local summary_prompt="以下の議論履歴を要約してください。各参加者の主張、合意点、未解決の論点を保持してください。

${CONVERSATION_HISTORY}

要約ルール:
- 各参加者の立場と根拠を正確に保持する
- 合意された事項は明確にマークする
- 未解決の論点を明示する
- 1000-2000文字程度に圧縮する"

    local summary
    if summary=$(call_haiku "$summary_prompt"); then
        CONVERSATION_HISTORY="[以下は議論のこれまでの要約です]

${summary}

"
        emit "SYSTEM" "コンテキストを圧縮しました"
    else
        emit "WARNING" "コンテキスト圧縮に失敗しました。圧縮なしで続行します。"
    fi
}

# =============================================================================
# 議論ループ
# =============================================================================

discussion_loop() {
    local round=1

    while [ "$round" -le "$MAX_ROUNDS" ]; do
        emit_header "ラウンド ${round} / ${MAX_ROUNDS}"

        # コンテキスト圧縮チェック
        if [ ${#CONVERSATION_HISTORY} -gt $MAX_CONTEXT_CHARS ]; then
            summarize_history
        fi

        # ファシリテーション
        local facilitation
        if ! facilitation=$(phase_facilitate "$round"); then
            return 1
        fi

        # 収束判定
        if check_convergence "$facilitation"; then
            emit_header "合意に到達しました (ラウンド ${round})"
            return 0
        fi

        # 応答
        phase_respond "$facilitation"

        round=$((round + 1))
    done

    emit_header "最大ラウンド数 (${MAX_ROUNDS}) に到達しました"
    return 0
}

# =============================================================================
# クリーンアップ
# =============================================================================

cleanup() {
    rm -f /tmp/discuss_opus_* /tmp/discuss_codex_* 2>/dev/null || true
}

trap cleanup EXIT
trap 'echo ""; emit "SYSTEM" "中断されました"; exit 130' INT

# =============================================================================
# メイン
# =============================================================================

main() {
    parse_args "$@"

    # 依存チェック
    for cmd in claude codex; do
        if ! command -v "$cmd" &> /dev/null; then
            echo "エラー: ${cmd} コマンドが見つかりません"
            exit 1
        fi
    done

    init_session
    collect_context
    build_system_prompts

    # Phase 1: 初期意見
    phase_initial_opinions

    # Phase 2-3: 議論ループ
    discussion_loop

    # Phase 4: 結論まとめ
    phase_conclude

    # ファイルに結果情報を追記
    {
        echo ""
        echo "---"
        echo ""
        echo "*議論完了: $(date '+%Y-%m-%d %H:%M:%S')*"
    } >> "$OUTPUT_FILE"

    emit_header "議論完了"
    emit "SYSTEM" "結果は ${OUTPUT_FILE} に保存されました"
}

main "$@"

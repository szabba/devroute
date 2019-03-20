#!/usr/bin/env bash

set -e

CONF_SCRIPT="$1"
source "$CONF_SCRIPT"
WORK_DIR="out"
RUNTIME="go111"
MEMORY="128M"

function ensure_workdir {
    mkdir -p out/
}

function sha256 {
    local of_file="$1"
    sha256sum "$of_file" | sed -e 's/\([0-9a-f]\{64\}\).*/\1/'
}

function pack {
    echo pack "$@" 1>&2
    local src_dir="$1"
    local final_segment
    final_segment=$(basename "$src_dir")
    local temporary_target
    temporary_target="$WORK_DIR/$final_segment.zip"
    zip -r -j "$temporary_target" "$src_dir"/ 1>&2
    local sum
    sum=$(sha256 "$temporary_target")
    local target="$WORK_DIR/$sum.zip"
    mv "$temporary_target" "$target" 1>&2
    echo "$target"
}

function upload_code {
    echo upload_code "$@" 1>&2
    local code_archive="$1"
    local archive_basename
    archive_basename="$(basename "\"$code_archive")"
    local target_url
    target_url="gs://$CODE_BUCKET/$archive_basename"
    gsutil cp -n "$code_archive" "$target_url" 1>&2
    echo "$target_url"
}

function deploy_uploaded_fn {
    echo deploy_uploaded_fn "$@" 1>&2
    local fn_name="$1"
    local entry_point="$2"
    local trigger="$3"
    local source="$4"
    gcloud functions deploy \
        --project "$PROJECT" \
        --region "$REGION" \
        --runtime "$RUNTIME" \
        --memory "$MEMORY" \
        "$fn_name" \
        --entry-point "$entry_point" \
        --source "$source" \
        "--trigger-$trigger"
}

function deploy_fn {
    echo deploy_fn "$@" 1>&2
    local src_dir="$1"
    local entry_point="$2"
    local trigger="$3"
    local fn_name
    fn_name=$(basename "$src_dir")
    local code_archive
    code_archive=$(pack "$src_dir")
    local code_url
    code_url=$(upload_code "$code_archive")
    deploy_uploaded_fn "$fn_name" "$entry_point" "$trigger" "$code_url" 1>&2
}

ensure_workdir
deploy_fn backend/functions/hello Hello http
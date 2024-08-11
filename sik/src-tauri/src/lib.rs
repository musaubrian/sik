use std::process::Command;

use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct SearchResult {
    full_path: String,
    filename: String,
    hits: Vec<String>,
}

// Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
#[tauri::command]
fn search(query: &str) -> String {
    let output = Command::new("/home/musaubrian/personal/sik/sik.py")
        .args(["--json", "-q", query])
        .output()
        .expect("Failed to run command <sik>");
    let json_str = String::from_utf8_lossy(&output.stdout).to_string();
    let results: Vec<SearchResult> =
        serde_json::from_str(&json_str).expect("Failed to deserialize json");

    serde_json::to_string(&results).expect("Failed to serialize results")
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .invoke_handler(tauri::generate_handler![search])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

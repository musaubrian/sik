use std::fs;
use std::process::Command;

use dirs::home_dir;
use open;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct SearchResult {
    full_path: String,
    filename: String,
    hits: Vec<String>,
}

#[derive(Debug, Deserialize)]
struct Sources {
    sik_source: String,
}

// Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
#[tauri::command]
fn search(query: &str) -> String {
    let user_home = match home_dir() {
        Some(path) => path,
        _none => {
            let err = SearchResult {
                full_path: "".to_string(),
                filename: "Error: Could not get home location".to_string(),
                hits: [].to_vec(),
            };
            return serde_json::to_string(&vec![err]).unwrap();
        }
    };

    let config_path = user_home.join(".sik").join("config.toml");
    let config: Sources = match fs::read_to_string(&config_path) {
        Ok(config_text) => match toml::from_str(&config_text) {
            Ok(config) => config,
            Err(_) => {
                let err = SearchResult {
                    full_path: "".to_string(),
                    filename: format!("Error: Could not parse config file at {:?}", config_path),
                    hits: [].to_vec(),
                };
                return serde_json::to_string(&vec![err]).unwrap();
            }
        },
        Err(_) => {
            let err = SearchResult {
                full_path: "".to_string(),
                filename: format!(
                    "Error: Could not read config file at {:?}\nHave you initialized sik",
                    config_path
                ),
                hits: [].to_vec(),
            };
            return serde_json::to_string(&vec![err]).unwrap();
        }
    };

    let output = match Command::new(&config.sik_source)
        .args(["--json", "-q", query])
        .output()
    {
        Ok(output) => output,
        Err(_) => {
            let err = SearchResult {
                full_path: "".to_string(),
                filename: format!(
                    "Error: Failed to run sik command from {:?}",
                    config.sik_source
                ),
                hits: [].to_vec(),
            };
            return serde_json::to_string(&vec![err]).unwrap();
        }
    };

    let json_str = String::from_utf8_lossy(&output.stdout).to_string();
    let results: Vec<SearchResult> = match serde_json::from_str(&json_str) {
        Ok(results) => results,
        Err(_) => {
            let err = SearchResult {
                full_path: "".to_string(),
                filename: "Error: Failed to parse sik output".to_string(),
                hits: [].to_vec(),
            };
            return serde_json::to_string(&vec![err]).unwrap();
        }
    };

    serde_json::to_string(&results).unwrap()
}

#[tauri::command]
fn open_file(url: &str) {
    open::that(url).unwrap()
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .invoke_handler(tauri::generate_handler![search, open_file])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

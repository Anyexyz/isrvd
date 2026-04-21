use poem::{handler, Route, Server, web::Json, listener::TcpListener, get, Result};
use serde::{Deserialize, Serialize};
use std::process::Command;
use std::env;

#[derive(Serialize, Deserialize)]
struct HealthCheck {
    status: String,
    message: String,
}

#[derive(Serialize, Deserialize)]
struct SystemInfo {
    os: String,
    hostname: String,
    poem_version: String,
}

#[derive(Serialize, Deserialize)]
struct IsrvdInfo {
    status: String,
    version: String,
    services: Vec<String>,
}

#[handler]
async fn health_check() -> Json<HealthCheck> {
    Json(HealthCheck {
        status: "ok".to_string(),
        message: "Poem service is running".to_string(),
    })
}

#[handler]
async fn get_system_info() -> Json<SystemInfo> {
    let os = env::consts::OS.to_string();
    let hostname = Command::new("hostname")
        .output()
        .map(|output| String::from_utf8_lossy(&output.stdout).trim().to_string())
        .unwrap_or_else(|_| "unknown".to_string());
    
    Json(SystemInfo {
        os,
        hostname,
        poem_version: "3.1.12".to_string(),
    })
}

#[handler]
async fn get_isrvd_info() -> Result<Json<IsrvdInfo>> {
    // 这里可以实现与 Isrvd Go 服务的交互
    // 暂时返回模拟数据
    Ok(Json(IsrvdInfo {
        status: "running".to_string(),
        version: "1.0.0".to_string(),
        services: vec!["docker", "swarm", "apisix", "compose"].iter().map(|s| s.to_string()).collect(),
    }))
}

#[tokio::main]
async fn main() -> Result<(), std::io::Error> {
    // 初始化日志
    tracing_subscriber::fmt::init();
    
    // 创建路由
    let app = Route::new()
        .at("/health", get(health_check))
        .at("/system", get(get_system_info))
        .at("/isrvd", get(get_isrvd_info));
    
    // 启动服务器
    Server::new(TcpListener::bind("127.0.0.1:3000"))
        .run(app)
        .await
}

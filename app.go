package main

import (
 "context"
 "encoding/json"
 "fmt"
 "io"
 "log"
 "net/http"
 "os"
 "os/signal"
 "syscall"
 "time"
)

const page = "<!doctype html>\n<html lang=\"ru\">\n<head>\n<meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\">\n<title>StudyOps — время учиться</title>\n<style>\n:root{color-scheme:dark;font-family:system-ui,sans-serif;background:#10131d;color:#eff3ff}*{box-sizing:border-box}body{margin:0;padding:32px 20px;background:radial-gradient(ellipse at 15% 0,#23374e 0,transparent 55%);min-height:100vh}main{max-width:960px;margin:auto}header,footer,.row{display:flex;justify-content:space-between;gap:16px;align-items:center;flex-wrap:wrap}.logo{font-weight:800;font-size:24px}.tag{font:12px monospace;color:#9ec9ff;border:1px solid #344157;padding:8px 12px;border-radius:30px}.eyebrow{margin-top:70px;letter-spacing:3px;color:#a0e8cc;font-size:12px}h1{font-size:clamp(36px,7vw,64px);max-width:700px;line-height:1.1;margin:18px 0}p{color:#b1bdd1;line-height:1.6}.grid{display:grid;grid-template-columns:1.5fr 1fr;gap:22px;margin:38px 0}.card{padding:30px;border:1px solid #313c50;border-radius:24px;background:#1a2130}.label{font-size:12px;text-transform:uppercase;letter-spacing:2px;color:#a9b7cd}.timer{font-size:clamp(60px,10vw,90px);font-variant-numeric:tabular-nums;font-weight:700;margin:20px 0}button,select{font:inherit;border-radius:12px;border:1px solid #45536a;background:#222f45;color:white;padding:12px 18px;cursor:pointer}button.primary{background:#a0e8cc;color:#102d27;border:0;font-weight:700}button:focus-visible,select:focus-visible{outline:3px solid #a0e8cc;outline-offset:4px}dl{margin:24px 0 0}dt{color:#b1bdd1;font-size:13px;margin-top:20px}dd{margin:7px 0;font-family:monospace;overflow-wrap:anywhere}footer{font-size:12px;color:#8998b1}.counter{color:#a0e8cc;margin-top:24px}.status{color:#a0e8cc}@media(max-width:680px){.grid{grid-template-columns:1fr}.eyebrow{margin-top:44px}.card{padding:24px}}\n</style>\n</head>\n<body><main>\n<header><div class=\"logo\">◈ StudyOps</div><span class=\"tag\">GO / GITHUB / LINUX</span></header>\n<div class=\"eyebrow\">МЕНЬШЕ ОТВЛЕЧЕНИЙ. БОЛЬШЕ РЕЗУЛЬТАТА.</div>\n<h1>Твой следующий<br>шаг — в фокусе.</h1>\n<p>Выбери время, запусти таймер и посвяти эту сессию одной задаче.</p>\n<div class=\"grid\"><section class=\"card\" aria-label=\"Таймер подготовки\">\n<div class=\"row\"><span class=\"label\">Фокус-сессия</span><select id=\"duration\" aria-label=\"Длительность сессии\"><option value=\"25\">25 минут</option><option value=\"15\">15 минут</option><option value=\"5\">5 минут</option></select></div>\n<div id=\"clock\" class=\"timer\" role=\"timer\" aria-label=\"Оставшееся время\">25:00</div>\n<div class=\"row\"><button id=\"start\" class=\"primary\">Начать сессию</button><button id=\"reset\">Сбросить</button></div>\n<div class=\"counter\" id=\"completed\"></div><p id=\"message\" role=\"status\">Готов к работе.</p>\n</section><aside class=\"card\"><div class=\"label\">Приложение в работе</div>\n<p class=\"status\" id=\"connection\" role=\"status\">Проверяю сервер…</p>\n<dl><dt>Окружение</dt><dd id=\"environment\">—</dd><dt>Версия (коммит)</dt><dd id=\"version\">—</dd><dt>Время работы сервера</dt><dd id=\"uptime\">—</dd></dl>\n<p>Версия помогает проверить, что после push на сервер попали изменения.</p>\n</aside></div><footer><span>Разработано на Go · StudyOps</span><span>GitLab Flow: feature → main → staging → production</span></footer>\n</main><script>\n'use strict';\nconst $=id=>document.getElementById(id);\nlet seconds=1500, deadline=0, running=false, completed=0;\ntry{completed=Number(localStorage.getItem('studyops.sessions'))||0;}catch{}\nfunction count(){$('completed').textContent=`Завершено сессий: ${completed}`;}count();\nfunction render(){$('clock').textContent=`${String(Math.floor(seconds/60)).padStart(2,'0')}:${String(seconds%60).padStart(2,'0')}`;document.title=`${$('clock').textContent} · StudyOps`;}\nfunction pause(){running=false;$('start').textContent='Продолжить';$('duration').disabled=false;}\nfunction tick(){if(!running)return;seconds=Math.max(0,Math.ceil((deadline-Date.now())/1000));render();if(seconds===0){pause();completed++;try{localStorage.setItem('studyops.sessions',String(completed));}catch{}count();$('start').textContent='Новая сессия';$('message').textContent='Сессия завершена! Сделай небольшой перерыв.';}}\n$('start').onclick=()=>{if(running){tick();if(!running)return;pause();$('message').textContent='Пауза.';return;}if(seconds===0)seconds=Number($('duration').value)*60;running=true;deadline=Date.now()+seconds*1000;$('start').textContent='Пауза';$('duration').disabled=true;$('message').textContent='Время сосредоточиться.';render();};\nfunction reset(){running=false;seconds=Number($('duration').value)*60;$('duration').disabled=false;$('start').textContent='Начать сессию';$('message').textContent='Готов к работе.';render();}\n$('reset').onclick=reset;$('duration').onchange=reset;setInterval(tick,250);\nasync function info(){try{const r=await fetch('/api/info',{cache:'no-store'});if(!r.ok)throw Error(r.status);const d=await r.json();$('environment').textContent=d.environment;$('version').textContent=d.version;$('uptime').textContent=`${Math.floor(d.uptime_seconds/60)} мин ${d.uptime_seconds%60} с`;$('connection').textContent='● Сервер доступен';}catch{$('connection').textContent='Сервер недоступен';}}info();setInterval(info,15000);\n</script></body></html>\n"

var version = "dev"
var started = time.Now()

func routes() http.Handler {
 mux := http.NewServeMux()
 mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.Header().Set("X-Content-Type-Options", "nosniff")
  _, _ = io.WriteString(w, page)
 })
 mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  _, _ = io.WriteString(w, "{\"status\":\"ok\"}\n")
 })
 mux.HandleFunc("GET /api/info", func(w http.ResponseWriter, r *http.Request) {
  env := os.Getenv("APP_ENV")
  if env == "" { env = "local" }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Cache-Control", "no-store")
  _ = json.NewEncoder(w).Encode(map[string]any{
   "version": version, "environment": env,
   "uptime_seconds": int(time.Since(started).Seconds()),
  })
 })
 return mux
}

func main() {
 if len(os.Args) == 2 && os.Args[1] == "--healthcheck" {
  client := &http.Client{Timeout: 2 * time.Second}
  resp, err := client.Get("http://127.0.0.1:8888/healthz")
  if err != nil { os.Exit(1) }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK { os.Exit(1) }
  return
 }
 server := &http.Server{Addr: ":8888", Handler: routes(), ReadHeaderTimeout: 5*time.Second, ReadTimeout: 10*time.Second, WriteTimeout: 10*time.Second, IdleTimeout: 60*time.Second}
 ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
 defer stop()
 go func() {
  <-ctx.Done()
  shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  _ = server.Shutdown(shutdown)
 }()
 fmt.Printf("StudyOps %s listening on :8888\n", version)
 if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal(err) }
}

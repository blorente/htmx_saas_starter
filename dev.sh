air -c air.toml &

AIR_PID=$!

cleanup() {
    kill -2 $AIR_PID || true
    main_pid=$(lsof -i tcp:9999 | tail -n 1 | awk '{print $2}')
    kill -2 "${main_pid}" || true
    exit 0
}

trap cleanup EXIT SIGINT SIGTERM

npx tailwindcss \
  -i 'styles.css' \
  -o 'public/styles.css' \
  --watch &

wait $AIR_PID

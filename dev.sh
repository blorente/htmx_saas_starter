clean_port() {
    main_pid=$(lsof -i tcp:9999 | tail -n 1 | awk '{print $2}')
    kill -2 "${main_pid}" || true
}

cleanup() {
    kill -2 $AIR_PID || true
    clean_port
    exit 0
}

air -c air.toml &
AIR_PID=$!
npx tailwindcss \
  -i 'styles.css' \
  -o 'public/styles.css' \
  --watch &

trap cleanup EXIT SIGINT SIGTERM

clean_port

wait $AIR_PID

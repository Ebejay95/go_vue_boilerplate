#!/bin/bash

echo "🔥 Starting Vue.js Development with Force Watch..."
echo "📁 Watching: /app/src for changes"

# Function to trigger rebuild
trigger_rebuild() {
    echo "🔄 File change detected, triggering rebuild..."
    # Touch a file to force webpack rebuild
    touch /app/src/.rebuild-trigger
    sleep 1
    rm -f /app/src/.rebuild-trigger
}

# Start Vue dev server in background
echo "🚀 Starting Vue dev server..."
npm run serve -- --host 0.0.0.0 --port ${FRONTEND_PORT} &
VUE_PID=$!

# Give Vue server time to start
sleep 5

echo "👀 Starting file watcher..."

# Watch for file changes and trigger rebuilds
while true; do
    find /app/src -name "*.vue" -o -name "*.js" -o -name "*.ts" -o -name "*.json" | while read file; do
        if [ -f "$file" ]; then
            # Get current modification time
            current_time=$(stat -c %Y "$file" 2>/dev/null || echo "0")

            # Check if we've seen this file before
            cache_file="/tmp/watch_$(echo "$file" | sed 's/[^a-zA-Z0-9]/_/g')"

            if [ -f "$cache_file" ]; then
                last_time=$(cat "$cache_file")
            else
                last_time=0
            fi

            # If file was modified, trigger rebuild
            if [ "$current_time" -gt "$last_time" ]; then
                echo "📝 Change detected in: $file"
                echo "$current_time" > "$cache_file"
                trigger_rebuild
            fi
        fi
    done

    # Check every 2 seconds
    sleep 2
done &
WATCH_PID=$!

echo "✅ Development server running with file watcher"
echo "🌐 Access your app at: http://localhost:${FRONTEND_PORT}"
echo "🔥 Hot reload is active - make changes to see live updates!"

# Handle shutdown gracefully
trap 'echo "🛑 Shutting down..."; kill $VUE_PID $WATCH_PID; exit 0' SIGTERM SIGINT

# Wait for Vue process
wait $VUE_PID

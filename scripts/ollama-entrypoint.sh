#!/bin/sh
# Ollama entrypoint with auto model pull

# 启动 ollama 服务（后台）
ollama serve &
OLLAMA_PID=$!

# 等待 ollama 就绪
echo "Waiting for Ollama to be ready..."
until curl -s http://localhost:11434/api/version > /dev/null 2>&1; do
    sleep 1
done
echo "Ollama is ready!"

# 拉取 embedding 模型
echo "Pulling embedding model: ${EMBEDDING_MODEL:-nomic-embed-text}"
ollama pull ${EMBEDDING_MODEL:-nomic-embed-text}

# 等待 ollama 进程
wait $OLLAMA_PID

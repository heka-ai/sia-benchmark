BENCHMARK_CMD="python3 benchmark_serving.py --backend openai \
    --base-url http://${EC2_PUBLIC_DNS_GPU}:8000 \
    --dataset-name ${DATASET_NAME} \
    --dataset-path ${DATASET_PATH} \
    --hf-revision ${HF_REVISION} \
    --hf-split ${HF_SPLIT} \
    --model ${MODEL_ID} \
    --seed 12345 \
    --save-result \
    --result-filename metrics.json \
    --num-prompts ${NUM_PROMPTS}"

echo "Running benchmark command:"
echo "$BENCHMARK_CMD"
eval "$BENCHMARK_CMD"
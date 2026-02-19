import json
import os
from datetime import datetime
from typing import Any, Dict, List


def save_result(
    deployment_config: dict,
    args,
    benchmark_result: dict,
    list_prompt: List[str],
    list_expected_output: List[str],
):
    """Save benchmark results to a JSON file with structured data."""
    result_json: Dict[str, Any] = {}
    environnement = {}
    models = {}
    results = {}
    tasks = {}
    datasets = {}
    benchmarks = {}
    evaluation = {}
    # Setup
    current_dt = datetime.now().strftime("%Y%m%d-%H%M%S")

    results["input"] = list_prompt
    results["expected_output"] = list_expected_output
    results["actual_output"] = benchmark_result["generated_texts"]
    results["itls"] = benchmark_result["itls"]
    results["ttfts"] = benchmark_result["ttfts"]
    results["input_lens"] = benchmark_result["input_lens"]
    results["output_lens"] = benchmark_result["output_lens"]

    provider = deployment_config.get("provider")
    environnement["provider"] = provider
    environnement["cpu_instance_type"] = deployment_config[provider].get("cpu_instance_type")
    environnement["gpu_instance_type"] = deployment_config[provider].get("gpu_instance_type")
    environnement["region"] = deployment_config[provider].get("region")

    models["name"] = args.model
    models["tokenizer_id"] = args.tokenizer
    models["url"] = args.model
    models["best_of"] = args.best_of

    datasets["url"] = args.dataset_path
    datasets["revision"] = args.hf_revision
    datasets["split"] = args.hf_split
    datasets["name"] = args.dataset_name

    metrics_dict = deployment_config["benchmark"].get("metricsdesired", {})
    benchmarks["date"] = current_dt
    benchmarks["metrics"] = metrics_dict
    benchmarks["request_rate"] = (
        args.request_rate if args.request_rate < float("inf") else "inf"
    )
    benchmarks["duration"] = benchmark_result["duration"]
    benchmarks["completed"] = benchmark_result["completed"]
    benchmarks["total_output_tokens"] = benchmark_result["total_output_tokens"]
    benchmarks["request_throughput"] = benchmark_result["request_throughput"]
    benchmarks["output_throughput"] = benchmark_result["output_throughput"]
    benchmarks["total_token_throughput"] = benchmark_result["total_token_throughput"]
    benchmarks["num_prompts"] = args.num_prompts

    # Evaluation parameters
    evaluation["evaluation_model"] = args.evaluation_model
    evaluation["prompt_template"] = args.prompt_template
    evaluation["top_k"] = args.top_k
    evaluation["metrics_desired"] = metrics_dict

    # Metadata
    if args.metadata:
        for item in args.metadata:
            if "=" in item:
                kvstring = item.split("=")
                result_json[kvstring[0].strip()] = kvstring[1].strip()
            else:
                raise ValueError(
                    "Invalid metadata format. Please use KEY=VALUE format."
                )

    result_json["results"] = results
    result_json["environment"] = environnement
    result_json["model"] = models
    result_json["task"] = tasks
    result_json["benchmark"] = benchmarks
    result_json["dataset"] = datasets
    result_json["metrics"] = metrics_dict
    result_json["evaluation"] = evaluation
    result_json["benchmark_id"] = deployment_config.get("benchmark_id")

    inference_engine = deployment_config.get("inference_engine")
    result_json["inference_engine"] = inference_engine
    result_json["deployment_params"] = deployment_config[inference_engine]
    result_json["deployment_params"] = {
        key: value for key, value in result_json["deployment_params"].items() if value not in ["", None, [], {}]
    }

    if args.result_filename:
        file_name = args.result_filename
    if args.result_dir:
        file_name = os.path.join(args.result_dir, file_name)
    with open(file_name, "w", encoding="utf-8") as outfile:
        json.dump(result_json, outfile)

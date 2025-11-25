import base64
import io
from typing import Collection, Dict, List, Optional, Tuple

import numpy as np
from datasets import load_dataset
from PIL import Image
from transformers import PreTrainedTokenizerBase


def sample_hf_requests(
    dataset_path: str,
    dataset_subset: str,
    dataset_split: str,
    num_requests: int,
    tokenizer: PreTrainedTokenizerBase,
    revision: str = "main",
    fixed_output_len: Optional[int] = None,
) -> List[Tuple[str, str, int, Optional[Dict[str, Collection[str]]]]]:
    dataset = load_dataset(
        dataset_path, name=dataset_subset, split=dataset_split, revision=revision
    )

    
    assert (
        "conversations" in dataset.features
    ), "HF Dataset must have 'conversations' column."
    filtered_dataset = dataset.shuffle().filter(lambda x: len(x["conversations"]) >= 2)
    sampled_requests: List[Tuple[str, int, int, Dict[str, Collection[str]]]] = []
    for data in filtered_dataset:
        if len(sampled_requests) == num_requests:
            break

        # Tokenize the prompts and completions.
        prompt = data["conversations"][0]["value"]
        prompt_token_ids = tokenizer(prompt).input_ids
        completion = data["conversations"][1]["value"]
        completion_token_ids = tokenizer(completion).input_ids
        prompt_len = len(prompt_token_ids)
        output_len = (
            len(completion_token_ids) if fixed_output_len is None else fixed_output_len
        )
        if fixed_output_len is None and (prompt_len < 4 or output_len < 4):
            # Prune too short sequences.
            continue
        if fixed_output_len is None and (
            prompt_len > 1024 or prompt_len + output_len > 2048
        ):
            # Prune too long sequences.
            continue

        if "image" in data and isinstance(data["image"], Image):
            image: Image = data["image"]
            image = image.convert("RGB")
            image_data = io.BytesIO()
            image.save(image_data, format="JPEG")
            image_base64 = base64.b64encode(image_data.getvalue()).decode("utf-8")
            mm_content = {
                "type": "image_url",
                "image_url": {"url": f"data:image/jpeg;base64,{image_base64}"},
            }
        else:
            mm_content = None

        sampled_requests.append((prompt,completion, prompt_len, output_len, mm_content))
    
    return sampled_requests


def sample_heka_requests(
    dataset_path: str,
    dataset_subset: str,
    dataset_split: str,
    num_requests: int,
    tokenizer: PreTrainedTokenizerBase,
    task: str,
    revision: str = "main",
    fixed_output_len: Optional[int] = None,
) -> List[Tuple[str, str, int, Optional[Dict[str, Collection[str]]]]]:
    dataset = load_dataset(
        dataset_path, name=dataset_subset, split=dataset_split, revision=revision
    )
    columns_required = ["input", "expected_output"]
    for column_name in columns_required:
        assert (
            column_name in dataset.features
        ), f"HF Dataset must have {column_name} column."

    filtered_dataset = dataset
    sampled_requests: List[Tuple[str, int, int, Dict[str, Collection[str]]]] = []
    list_prompt = []
    list_expected_output = []
    for data in filtered_dataset:
        if len(sampled_requests) == num_requests:
            break
        # Tokenize the prompts and completions.
        prompt = data["input"]
        prompt_token_ids = tokenizer(prompt).input_ids
        completion = data["expected_output"]
        completion_token_ids = tokenizer(completion).input_ids
        prompt_len = len(prompt_token_ids)
        output_len = (
            len(completion_token_ids) if fixed_output_len is None else fixed_output_len
        )
        if fixed_output_len is None and (prompt_len < 4 or output_len < 4):
            # Prune too short sequences.
            continue
        if fixed_output_len is None and (
            prompt_len > 1024 or prompt_len + output_len > 2048
        ):
            # Prune too long sequences.
            continue

        if "image" in data and isinstance(data["image"], Image):
            image: Image = data["image"]
            image = image.convert("RGB")
            image_data = io.BytesIO()
            image.save(image_data, format="JPEG")
            image_base64 = base64.b64encode(image_data.getvalue()).decode("utf-8")
            mm_content = {
                "type": "image_url",
                "image_url": {"url": f"data:image/jpeg;base64,{image_base64}"},
            }
        else:
            mm_content = None

        sampled_requests.append((prompt, prompt_len, output_len, mm_content))
        list_prompt.append(prompt)
        list_expected_output.append(completion)
    
    return sampled_requests, list_prompt, list_expected_output


def sample_random_requests(
    prefix_len: int,
    input_len: int,
    output_len: int,
    num_prompts: int,
    range_ratio: float,
    tokenizer: PreTrainedTokenizerBase,
) -> List[Tuple[str, int, int]]:
    prefix_token_ids = np.random.randint(
        0, tokenizer.vocab_size, size=prefix_len
    ).tolist()

    input_lens = np.random.randint(
        int(input_len * range_ratio),
        input_len + 1,
        size=num_prompts,
    )
    output_lens = np.random.randint(
        int(output_len * range_ratio),
        output_len + 1,
        size=num_prompts,
    )
    offsets = np.random.randint(0, tokenizer.vocab_size, size=num_prompts)
    input_requests = []
    for i in range(num_prompts):
        prompt = tokenizer.decode(
            prefix_token_ids
            + [
                (offsets[i] + i + j) % tokenizer.vocab_size
                for j in range(input_lens[i])
            ]
        )

        input_requests.append(
            (prompt, int(prefix_len + input_lens[i]), int(output_lens[i]), None)
        )

    return input_requests

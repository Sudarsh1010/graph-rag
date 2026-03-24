import json
import glob
import sys
from sentence_transformers import SentenceTransformer

DATASET_DIR = "sap-order-to-cash-dataset"
OUTPUT_DIR = "data/embeddings"
MODEL_NAME = "all-MiniLM-L6-v2"


def generate_embeddings(entity_dir, output_path, text_field, id_fields):
    texts = []
    ids = []

    files = sorted(glob.glob(f"{DATASET_DIR}/{entity_dir}/part-*.jsonl"))
    for file in files:
        with open(file) as f:
            for line in f:
                record = json.loads(line)
                text = record.get(text_field, "")
                if not text or not text.strip():
                    continue

                id_key = "_".join(str(record.get(k, "")) for k in id_fields)
                if not id_key.strip("_"):
                    continue

                texts.append(text)
                ids.append(id_key)

    if not texts:
        print(f"  No text records found in {entity_dir}")
        return

    print(f"  Generating embeddings for {entity_dir}: {len(texts)} records")
    embeddings = model.encode(texts, show_progress_bar=True, batch_size=128)

    with open(output_path, "w") as f:
        for id_, emb in zip(ids, embeddings):
            f.write(json.dumps({"id": id_, "embedding": emb.tolist()}) + "\n")

    print(f"  Saved to {output_path}")


if __name__ == "__main__":
    print(f"Loading model: {MODEL_NAME}")
    model = SentenceTransformer(MODEL_NAME)

    generate_embeddings(
        "product_descriptions",
        f"{OUTPUT_DIR}/product_descriptions_embeddings.jsonl",
        "productDescription",
        ["product", "language"],
    )

    generate_embeddings(
        "business_partners",
        f"{OUTPUT_DIR}/business_partners_embeddings.jsonl",
        "businessPartnerFullName",
        ["businessPartner"],
    )

    generate_embeddings(
        "plants",
        f"{OUTPUT_DIR}/plants_embeddings.jsonl",
        "plantName",
        ["plant"],
    )

    print("Done!")

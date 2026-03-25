ALTER TABLE product_descriptions ADD COLUMN IF NOT EXISTS embedding vector(384);
ALTER TABLE business_partners ADD COLUMN IF NOT EXISTS embedding vector(384);
ALTER TABLE plants ADD COLUMN IF NOT EXISTS embedding vector(384);

CREATE INDEX IF NOT EXISTS idx_product_descriptions_embedding
    ON product_descriptions USING hnsw (embedding vector_cosine_ops);
CREATE INDEX IF NOT EXISTS idx_business_partners_embedding
    ON business_partners USING hnsw (embedding vector_cosine_ops);
CREATE INDEX IF NOT EXISTS idx_plants_embedding
    ON plants USING hnsw (embedding vector_cosine_ops);

these how you can count the TF-IDF matrices :
TF(word, doc) = freq(word, doc) / max_freq(doc)
IDF(word) = log2( total_docs / docs_containing_word )

examples :
Docs = 3

"go" appears in 2 docs

IDF = log2(3 / 2) ≈ 0.585 (num of docs / appearance)
TF: doc1: freq=3, max=6 → 0.5 (3/6)
doc2: freq=2, max=4 → 0.5
TF-IDF: 0.5 \* 0.585 ≈ 0.292

---

these when finishing the project :
⚠️ Next critical step (don’t skip)

Before indexing, you still need:

👉 URL normalization (avoid duplicates like / vs /index)
👉 domain restriction
👉 storage layer

⚠️ Important note

Right now:

❌ duplicates across runs (no persistence tracking)
❌ no compression
❌ no metadata (date, headers)

👉 totally fine for your stage

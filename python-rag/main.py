from fastapi import FastAPI
from pydantic import BaseModel
from utils import build_index, get_llm, query_policy

app = FastAPI()
index = build_index()
llm = get_llm()

class Query(BaseModel):
    query: str

@app.post("/query")
def query(q: Query):
    return {"response": query_policy(index, llm, q.query)}
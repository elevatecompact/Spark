from pydantic import BaseModel


class TransactionData(BaseModel):
    amount: float
    currency: str
    payment_method: str
    ip_address: str | None = None
    device_fingerprint: str | None = None
    timestamp: str | None = None


class UserContext(BaseModel):
    user_id: str
    account_age_days: float | None = None
    previous_transactions: int | None = None
    email_verified: bool | None = None
    phone_verified: bool | None = None


class FraudScoreRequest(BaseModel):
    transaction: TransactionData
    user: UserContext | None = None


class FraudScoreResponse(BaseModel):
    score: float
    risk_level: str
    flags: list[str]
    decision: str
    model_version: str = "noop-v1"


class FraudRule(BaseModel):
    rule_id: str
    name: str
    description: str
    severity: str
    category: str


class FraudAnalysisResponse(BaseModel):
    rules_triggered: list[FraudRule]
    overall_score: float
    risk_level: str
    decision: str


class BatchFraudRequest(BaseModel):
    requests: list[FraudScoreRequest]


class BatchFraudResponse(BaseModel):
    results: list[FraudScoreResponse]


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "fraud-ai"
    version: str = "0.1.0"

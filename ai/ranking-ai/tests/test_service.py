import pytest

from src.models import (
    RankContext,
    RankItem,
    RankingRequest,
    TrainingData,
)
from src.repository import RankingRepository
from src.events import NoopProducer
from src.service import RankingService


@pytest.fixture
def service():
    repo = RankingRepository()
    event_pub = NoopProducer()
    return RankingService(repo=repo, event_pub=event_pub)


@pytest.mark.asyncio
async def test_rank_items_returns_sorted_results(service):
    request = RankingRequest(
        items=[
            RankItem(item_id="a", features={"f1": 1.0, "f2": 2.0}),
            RankItem(item_id="b", features={"f1": 5.0}),
            RankItem(item_id="c", features={"f1": 0.0, "f2": 0.0}),
        ],
        context=RankContext(user_id="u1"),
    )

    response = await service.rank_items(request)

    assert len(response.ranked_items) == 3
    assert response.ranked_items[0].item_id == "b"
    assert response.ranked_items[0].score == 5.0
    assert response.ranked_items[0].position == 0
    assert response.ranked_items[1].item_id == "a"
    assert response.ranked_items[1].score == 3.0
    assert response.ranked_items[1].position == 1
    assert response.ranked_items[2].item_id == "c"
    assert response.ranked_items[2].score == 0.0
    assert response.ranked_items[2].position == 2
    assert response.model_version == "noop-v1"


@pytest.mark.asyncio
async def test_train_model_returns_expected_structure(service):
    data = [
        TrainingData(item_id="a", features={"f1": 1.0}, label=0.5),
        TrainingData(item_id="b", features={"f1": 2.0}, label=1.0),
    ]

    response = await service.train_model(data)

    assert response.status == "ok"
    assert response.samples_trained == 2
    assert response.model_version == "noop-v1"


@pytest.mark.asyncio
async def test_batch_rank_processes_multiple(service):
    req1 = RankingRequest(
        items=[RankItem(item_id="a", features={"f1": 1.0})]
    )
    req2 = RankingRequest(
        items=[
            RankItem(item_id="x", features={"f1": 5.0}),
            RankItem(item_id="y", features={"f1": 2.0}),
        ]
    )

    results = await service.batch_rank([req1, req2])

    assert len(results) == 2
    assert len(results[0].ranked_items) == 1
    assert results[0].ranked_items[0].item_id == "a"
    assert len(results[1].ranked_items) == 2
    assert results[1].ranked_items[0].item_id == "x"
    assert results[1].ranked_items[1].item_id == "y"

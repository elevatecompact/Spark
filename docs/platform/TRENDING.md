# Trending Algorithms

The trending system identifies and surfaces content that is gaining popularity rapidly, helping viewers discover what is currently popular on the SPARK platform.

## Trending Score Calculation

The trending score is calculated using a combination of metrics weighted by time decay. View velocity measures the rate of views per hour normalized against the content's typical performance. Engagement rate considers likes, comments, shares, and saves as a percentage of views. Creator boost amplifies trending signals for newer creators to support discovery. Content freshness gives higher weight to recent engagement, with a half-life of 6 hours.

## Trending Categories

Trending content is categorized into several surfaces for different discovery contexts. Global trending shows the most popular content across the entire platform. Category trending surfaces trending content within specific categories. Local trending shows trending content based on the user's geographic region. Social trending surfaces content trending within the user's social graph. Creator trending identifies creators who are rapidly gaining followers and engagement.

## Time Windows

Trending is calculated over multiple time windows to capture different patterns. Hot content shows what is trending in the last 2 hours for real-time discovery. Rising content shows what has been gaining traction over the last 24 hours. Weekly trending aggregates engagement over 7 days for a broader view.

## Abuse Prevention

The trending system includes safeguards against manipulation. Velocity caps prevent artificial inflation from bot traffic. Engagement quality scoring filters out low-quality engagement signals. Creator diversity ensures no single creator dominates trending categories. Reporting and review mechanisms allow manual intervention for suspected manipulation.

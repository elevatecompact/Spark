# Blue-Green Deployment

## Strategy

Titan uses a blue-green deployment model for production, maintaining two identical environments (blue and green) behind a load balancer. At any time, only one environment receives live traffic.

## How It Works

The active environment (e.g., blue) serves all user traffic. The inactive environment (green) is deployed with the new version. Automated smoke tests and health checks validate the green environment. The load balancer then switches from blue to green. Green becomes active; blue becomes standby.

## Implementation

In Kubernetes, two sets of Deployments and Services are labeled `deployment=blue` and `deployment=green`. Traffic switch is achieved via a single CNAME or ALB target group change through Terraform. The service mesh (Istio) VirtualService routes 100% of traffic to the active label.

## Automated Verification

Before traffic switch, the CD pipeline validates that all pods are Ready and passing probes, the synthetic health endpoint returns 200 with correct version metadata, integration tests pass against the green environment, and metrics (error rate, latency) are within acceptable bounds.

## Rollback

If metrics degrade within 15 minutes of switchover, the load balancer points back to the previous environment, a PagerDuty incident is automatically created, and the deployment is marked as failed.

## Benefits

Zero-downtime deployments, rapid rollback (DNS TTL or load balancer switch within seconds), and full isolation between versions eliminates traffic mixing during deployment. Database migrations are handled reversibly with backward-compatible schema changes only.
$headers = @{ Accept = "application/vnd.github+json" }
$base = "https://api.github.com/repos/elevatecompact/Spark"

$runs = Invoke-RestMethod -Uri "$base/actions/runs?per_page=1&status=failure" -Headers $headers
$run = $runs.workflow_runs[0]

$jobs = Invoke-RestMethod -Uri "$base/actions/runs/$($run.id)/jobs" -Headers $headers
foreach ($job in $jobs.jobs) {
    if ($job.conclusion -eq "failure") {
        Write-Host "Job: $($job.name) (id: $($job.id))"
        # List steps and their conclusions
        foreach ($step in $job.steps) {
            Write-Host "  Step: $($step.name) -> $($step.conclusion) ($($step.started_at))"
        }
    }
}

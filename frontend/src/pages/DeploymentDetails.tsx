import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { deploymentsApi } from '@/lib/api';
import type { Deployment, DeploymentLogs } from '@/lib/types';
import { extractString } from '@/lib/types';
import StatusBadge from '@/components/StatusBadge';
import LogsViewer from '@/components/LogsViewer';

export default function DeploymentDetailsPage() {
  const params = useParams<{ id: string; deploymentId: string }>();
  const deploymentId = params.deploymentId!;
  const appId = params.id!;

  const [deployment, setDeployment] = useState<Deployment | null>(null);
  const [logs, setLogs] = useState<DeploymentLogs | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (deploymentId) {
      loadDeployment();
      loadLogs();
    }
  }, [deploymentId]);

  const loadDeployment = async () => {
    try {
      setError(null);
      const data = await deploymentsApi.getById(deploymentId);
      setDeployment(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load deployment');
      console.error('Error loading deployment:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadLogs = async () => {
    try {
      const data = await deploymentsApi.getLogs(deploymentId);
      setLogs(data);
    } catch (err) {
      console.error('Error loading logs:', err);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-[var(--app-bg)] flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--primary)]"></div>
          <p className="mt-4 text-[var(--text-secondary)]">Loading deployment...</p>
        </div>
      </div>
    );
  }

  if (error || !deployment) {
    return (
      <div className="min-h-screen bg-[var(--app-bg)]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Link
            to={`/apps/${appId}`}
            className="text-[var(--info)] hover:text-[var(--primary)] mb-6 inline-block transition-colors"
          >
            ← Back to App
          </Link>
          <div className="bg-[var(--error)]/10 border border-[var(--error)] rounded-lg p-6">
            <p className="text-[var(--error)]">{error || 'Deployment not found'}</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[var(--app-bg)]">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Link
          to={`/apps/${appId}`}
          className="text-[var(--info)] hover:text-[var(--primary)] mb-6 inline-block transition-colors"
        >
          ← Back to App
        </Link>

        <div className="bg-[var(--surface)] rounded-lg border border-[var(--border-subtle)] p-8 mb-6">
          <div className="flex items-start justify-between mb-6">
            <div>
              <h1 className="text-3xl font-bold text-[var(--text-primary)] mb-2">
                Deployment #{deployment.id}
              </h1>
              <StatusBadge status={deployment.status} />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div>
              <h3 className="text-sm font-medium text-[var(--text-muted)] mb-1">App ID</h3>
              <p className="text-[var(--text-primary)]">{deployment.app_id}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-[var(--text-muted)] mb-1">Status</h3>
              <StatusBadge status={deployment.status} />
            </div>
            {extractString(deployment.subdomain) && (
              <div>
                <h3 className="text-sm font-medium text-[var(--text-muted)] mb-1">Subdomain</h3>
                <p className="text-[var(--text-primary)]">{extractString(deployment.subdomain)}</p>
              </div>
            )}
            {extractString(deployment.image_name) && (
              <div>
                <h3 className="text-sm font-medium text-[var(--text-muted)] mb-1">Image Name</h3>
                <p className="text-[var(--text-primary)] font-mono text-sm">{extractString(deployment.image_name)}</p>
              </div>
            )}
            {extractString(deployment.container_id) && (
              <div>
                <h3 className="text-sm font-medium text-[var(--text-muted)] mb-1">Container ID</h3>
                <p className="text-[var(--text-primary)] font-mono text-sm">{extractString(deployment.container_id)}</p>
              </div>
            )}
            <div>
              <h3 className="text-sm font-medium text-[var(--text-muted)] mb-1">Created</h3>
              <p className="text-[var(--text-primary)]">
                {new Date(deployment.created_at).toLocaleString()}
              </p>
            </div>

            <div>
              <h3 className="text-sm font-medium text-[var(--text-muted)] mb-1">Updated</h3>
              <p className="text-[var(--text-primary)]">
                {new Date(deployment.updated_at).toLocaleString()}
              </p>
            </div>
          </div>

          {logs?.error_message && (
            <div className="mt-6 p-4 bg-[var(--error)]/10 border border-[var(--error)] rounded-lg">
              <h3 className="text-sm font-medium text-[var(--error)] mb-2">Error Message</h3>
              <p className="text-[var(--error)]">{extractString(logs.error_message)}</p>
            </div>
          )}
        </div>

        <div className="space-y-6">
          {logs && extractString(logs.build_log) && (
            <div>
              <h2 className="text-xl font-bold text-[var(--text-primary)] mb-4">Build Logs</h2>
              <LogsViewer logs={extractString(logs.build_log)} title="Build Logs" />
            </div>
          )}

          {logs && !extractString(logs.build_log) && (
            <div className="bg-[var(--surface)] rounded-lg border border-[var(--border-subtle)] p-8 text-center">
              <p className="text-[var(--text-secondary)]">No build logs available yet</p>
            </div>
          )}

          {logs && extractString(logs.runtime_log) && (
            <div>
              <h2 className="text-xl font-bold text-[var(--text-primary)] mb-4">Runtime Logs</h2>
              <LogsViewer logs={extractString(logs.runtime_log)} title="Runtime Logs" />
            </div>
          )}

          {logs && !extractString(logs.runtime_log) && (
            <div className="bg-[var(--surface)] rounded-lg border border-[var(--border-subtle)] p-8 text-center">
              <p className="text-[var(--text-secondary)]">No runtime logs available yet</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}



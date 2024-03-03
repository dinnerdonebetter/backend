import opentelemetry, { Tracer } from '@opentelemetry/api';
import { AlwaysOnSampler, NodeTracerProvider } from '@opentelemetry/sdk-trace-node';
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { TraceExporter } from '@google-cloud/opentelemetry-cloud-trace-exporter';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { Resource } from '@opentelemetry/resources';

// Enable OpenTelemetry exporters to export traces to Google Cloud Trace.
// Exporters use Application Default Credentials (ADCs) to authenticate.
// See https://developers.google.com/identity/protocols/application-default-credentials
// for more details.
const provider = new NodeTracerProvider({
  sampler: new AlwaysOnSampler(),
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: 'ddb-admin-server',
    [SemanticResourceAttributes.SERVICE_VERSION]: '1.0.0',
  }),
});

// Initialize the exporter. When your application is running on Google Cloud,
// you don't need to provide auth credentials or a project id.
const exporter = new TraceExporter();

// Configure the span processor to send spans to the exporter
provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

opentelemetry.trace.setGlobalTracerProvider(provider);

export const serverSideTracer: Tracer = opentelemetry.trace.getTracer('web-app-server');

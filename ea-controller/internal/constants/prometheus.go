package constants

// --- METRICS ---

const (
        MetricIstioRequests = "istio_requests_total"
        MetricKeplerEnergy  = "kepler_container_joules_total"
)

// --- LABELS ---

const (
        LabelServiceName      = "destination_service_name"
        LabelServiceNamespace = "destination_service_namespace"
        LabelReporter         = "reporter"
        LabelWorkload         = "destination_workload"
        LabelPod              = "pod_name"
        LabelContainerName    = "container_name"
        LabelContainerNS      = "container_namespace"
        LabelMode             = "mode"
)

// --- LABEL VALUES ---

const (
        ReporterDestination = "destination"

        ContainerIstioProxy = "istio-proxy"
        ContainerIstioInit  = "istio-init"

        ModeDynamic = "dynamic"
        ModeIdle    = "idle"
)

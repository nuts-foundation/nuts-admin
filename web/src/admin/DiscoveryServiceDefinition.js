export default class DiscoveryServiceDefinition {
    constructor(input) {
        // copy all properties to the class
        Object.assign(this, input);
    }

    // credentials returns an array of required credentials for this service.
    // Each entry is an array which contains the required types of credentials.
    credentials() {
        const pd = this.presentation_definition;
        return pd.input_descriptors.map((inputDescriptor) => {
            return inputDescriptor.constraints.fields.filter(f => f.path.includes('$.type') && f.filter && f.filter.type === "string")
                .map(constraint => constraint.filter.const)
        }).flat().filter(t => t !== "VerifiableCredential");
    }
}
export default {
    fields: [
        {
            name: "Name",
            description: "The name of the organization",
        },
        {
            name: "Location",
            description: "The name of the city or municipality where the organization is located",
        },
    ],

    render: (issuerDID, subjectDID, fieldValues) => {
        return {
            "@context": [
                "https://nuts.nl/credentials/v1",
                "https://www.w3.org/2018/credentials/v1",
            ],
            "issuer": issuerDID,
            "credentialSubject": {
                "id": subjectDID,
                "organization": {
                    "name": fieldValues[0],
                    "city": fieldValues[1],
                }
            },
            "type": [
                "NutsOrganizationCredential",
                "VerifiableCredential"
            ],
        }
    }
}
export default {
    fields: [
        {
            name: "URA",
            description: "UZI register abonneenummer (URA)",
        },
        {
            name: "Name",
            description: "Name of the care organization",
        },
        {
            name: "City",
            description: "Location where the care organization is based",
        },
    ],

    render: (issuerDID, subjectDID, fieldValues) => {
        return {
            "@context": [
                "https://nuts.nl/credentials/2024",
                "https://www.w3.org/2018/credentials/v1",
            ],
            "issuer": issuerDID,
            "credentialSubject": {
                "id": subjectDID,
                "organization": {
                    "ura": fieldValues[0],
                    "name": fieldValues[1],
                    "city": fieldValues[2],
                }
            },
            "type": [
                "NutsUraCredential",
                "VerifiableCredential"
            ],
        }
    }
}
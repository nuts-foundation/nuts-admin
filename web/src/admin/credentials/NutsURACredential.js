export default {
    fields: [
        {
            name: "URA",
            description: "UZI register abonneenummer (URA)",
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
                    "ura": fieldValues[0]
                }
            },
            "type": [
                "NutsURACredential",
                "VerifiableCredential"
            ],
        }
    }
}
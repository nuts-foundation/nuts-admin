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
                "https://nuts-services.nl/jsonld/credentials/experimental",
                "https://www.w3.org/2018/credentials/v1",
            ],
            "issuer": issuerDID,
            "credentialSubject": {
                "id": subjectDID,
                "ura": fieldValues[0]
            },
            "type": [
                "URACredential",
                "VerifiableCredential"
            ],
        }
    }
}
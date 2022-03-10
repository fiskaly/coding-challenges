# Signature Service - Coding Challenge

## Instructions

This challenge is part of the SWE interview process at fiskaly.

If you see this challenge, you've passed the first round of interviews and are now at the second and last stage.

We would ask you to prepare the challenge below to have a basis for discussion in the skill fit interview with two colleagues from the development department.

`quality > quantity`

### Project Setup

In order for you to not start from scratch with everything we provide you with:

- Go project containing setup
- Basic API structure and functionality
- Encoding / decoding of different key types (only needed when you want to serialize keys to a persistent storage)
- Key generation (ECC, RSA)
- Library to generate UUIDs included in `go.mod`

You can use these things as a foundation, but you're also free to modify them as you see fit.

### Prerequisites & Tooling

- Golang (v1.16+)

### The Challenge

The aim is to implement an API service that allows customers to create "signature devices" with which they can sign arbitrary transaction data.

#### Domain Description

The signature service can manage multiple signature devices. Such a device is identified by a unique identifier (e.g. UUID). For now you can pretend there is only one user / organisation using the system (e.g. a dedicated node for them) therefore you do not need to think about user management at all.

When creating the signature device, the client of the API has to make a choice and specify the signature algorithm that the device will be using to sign transaction data. During the creation process, a new key pair (public key & private key) has to be generated and assigned to the device.

The signature device should also have a `label` that can be used to display it in a UI and a `signature_counter` that tracks how many signatures have been performed with this device. While the `label` is provided by the user, the `signature_counter` shall only be modified internally.

##### Signature Creation

For the signature creation the client will have to provide `data_to_be_signed` through the API. In order to increase the security of the system, we will extend this raw data with the current `signature_counter` as well as the `last_signature`.

The resulting string should follow this format: `<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>`

For the base case that there is no `last_signature` (= `signature_counter == 0`) we will use the `base64` encoded device ID (`last_signature = base64(device.id)`).

This special string will be signed (`Signer.sign(secured_data_to_be_signed)`) and the resulting signature (`base64` encoded) will be returned to the client. A signature response could look like this:

```json
{ 
    "signature": <signature_base64_encoded>,
    "signed_data": "<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>"
}
```

After the signature has been created, the signature counters value shall be incremented (`signature_counter += 1`).

#### API

For now we need to provide two operations to our customers:

- `CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse`
- `SignTransaction(deviceId: string, data: string): SignatureResponse`

Think of how to expose these operations through a RESTful HTTP based API.

List / retrieval operations can optionally be implemented but aren't necessary by any means.

#### QA / Testing

As we're in the area of compliance technology we need to make sure that our implementation is verifiably correct. Think of an automatable way to assure the correctness (in this challenge: adherence to the specifications) of the system.

#### Technical Constraints & Considerations

- The system will be used by many concurrent clients accessing the same resources.
- The `signature_counter` shall be strictly monotonically increasing and ideally without any gaps.
- The system currently only supports `RSA` and `ECDSA` as signature algorithms. Try to design the signing mechanism in a way that allows easy extension to other algorithms without changing the core domain logic.
- For now it is enough to store signature devices in memory. Efficiency is not a priority for this. We might want to scale out though, therefore keep in mind that we might want to switch to a relational database at some point when designing your storage logic.

#### Credits

This challenge is heavily influenced by the KassenSichV (Germany) as well as the RKSV (Austria) and our solutions for them.

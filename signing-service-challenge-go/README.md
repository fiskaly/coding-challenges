# Signature Service - Coding Challenge

## Instructions

This challenge is part of the software engineering interview process at fiskaly.

If you see this challenge, you've passed the first round of interviews and are now at the second and last stage.

We would like you to attempt the challenge below. You will then be able to discuss your solution in the skill-fit interview with two of our colleagues from the development department.

The quality of your code is more important to us than the quantity.

### Project Setup

For the challenge, we provide you with:

- Go project containing the setup
- Basic API structure and functionality
- Encoding / decoding of different key types (only needed to serialize keys to a persistent storage)
- Key generation algorithms (ECC, RSA)
- Library to generate UUIDs, included in `go.mod`

You can use these things as a foundation, but you're also free to modify them as you see fit.

### Prerequisites & Tooling

- Golang (v1.20+)

### The Challenge

The goal is to implement an API service that allows customers to create `signature devices` with which they can sign arbitrary transaction data.

#### Domain Description

The `signature service` can manage multiple `signature devices`. Such a device is identified by a unique identifier (e.g. UUID). For now you can pretend there is only one user / organization using the system (e.g. a dedicated node for them), therefore you do not need to think about user management at all.

When creating the `signature device`, the client of the API has to choose the signature algorithm that the device will be using to sign transaction data. During the creation process, a new key pair (`public key` & `private key`) has to be generated and assigned to the device.

The `signature device` should also have a `label` that can be used to display it in the UI and a `signature_counter` that tracks how many signatures have been created with this device. The `label` is provided by the user. The `signature_counter` shall only be modified internally.

##### Signature Creation

For the signature creation, the client will have to provide `data_to_be_signed` through the API. In order to increase the security of the system, we will extend this raw data with the current `signature_counter` and the `last_signature`.

The resulting string (`secured_data_to_be_signed`) should follow this format: `<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>`

In the base case there is no `last_signature` (= `signature_counter == 0`). Use the `base64`-encoded device ID (`last_signature = base64(device.id)`) instead of the `last_signature`.

This special string will be signed (`Signer.sign(secured_data_to_be_signed)`) and the resulting signature (`base64` encoded) will be returned to the client. The signature response could look like this:

```json
{ 
    "signature": <signature_base64_encoded>,
    "signed_data": "<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>"
}
```

After the signature has been created, the signature counter's value has to be incremented (`signature_counter += 1`).

#### API

For now we need to provide two main operations to our customers:

- `CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse`
- `SignTransaction(deviceId: string, data: string): SignatureResponse`

Think of how to expose these operations through a RESTful HTTP-based API.

In addition, `list / retrieval operations` for the resources generated in the previous operations should be made available to the customers.

#### QA / Testing

As we are in the business of compliance technology, we need to make sure that our implementation is verifiably correct. Think of an automatable way to assure the correctness (in this challenge: adherence to the specifications) of the system.

#### Technical Constraints & Considerations

- The system will be used by many concurrent clients accessing the same resources.
- The `signature_counter` has to be strictly monotonically increasing and ideally without any gaps.
- The system currently only supports `RSA` and `ECDSA` as signature algorithms. Try to design the signing mechanism in a way that allows easy extension to other algorithms without changing the core domain logic.
- For now it is enough to store signature devices in memory. Efficiency is not a priority for this. In the future we might want to scale out. As you design your storage logic, keep in mind that we may later want to switch to a relational database.

#### Credits

This challenge is heavily influenced by the regulations for `KassenSichV` (Germany) as well as the `RKSV` (Austria) and our solutions for them.

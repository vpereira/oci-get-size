# oci-get-size

A tool to retrieve the uncompressed size of OCI images, taking into account multiple architectures.

## Features

- Retrieve uncompressed size of OCI images.
- Support for multiple architectures.

## Getting Started

These instructions will help you set up and run the project on your local machine.

### Prerequisites

Make sure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.16 or later)

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/oci-get-size.git
   cd oci-get-size
   ```

2. Build the project:
   ```sh
   make build
   ```

### Running the Project

After building the project, you can run it using:
```sh
./oci-get-size
```

### Submitting Images

To retrieve the uncompressed size of an image, submit it like this:
```sh
curl http://localhost:8080/get-uncompressed-size?image=registry.suse.com/bci/bci-busybox:latest
```

Example output:
```sh
$ curl http://localhost:8080/get-uncompressed-size?image=registry.suse.com/bci/bci-busybox:latest
{"image":"registry.suse.com/bci/bci-busybox:latest","sizes":{"amd64":13494784,"arm64":12823040,"ppc64le":15032832,"s390x":12002304}}
$
```

### Running Tests

To run the tests, use:
```sh
make test
```

### Formatting Code

To format the code, use:
```sh
make format
```

## Usage

Provide examples on how to use the project. For example:
```sh
./oci-get-size --option value
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

P4p170 (aka @josegomezr) for the ideas and debug sessions

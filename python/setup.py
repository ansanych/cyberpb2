from setuptools import setup, find_packages

setup(
    name="cyberpb2",
    version="0.1.1",
    description="Generated protobuf/gRPC package for Cybertele services",
    packages=find_packages(include=["cyberpb2", "cyberpb2.*"]),
    python_requires=">=3.10",
    install_requires=[
        "grpcio>=1.60",
        "protobuf>=5.0",
    ],
)
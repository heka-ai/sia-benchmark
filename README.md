# Sia Benchmark

The Sia Benchmark is a generalist LLM benchmark made to compare LLM on multiple axes :

- Performances
- Ecology
- Cost
- Intelligence

The benchmark tool is run and maintained by [Sia](https://sia-partners.com) employees. The benchmark results are published regularly on the [Sia Benchmark website](#).

### What does the benchmark tool do ?

The CLI is made to run a complete benchmark on any LLM, on most of the popular cloud providers. It will :

1. Create the instances on the cloud provider
2. Check the connection to the instances
3. Deploy the model on the instances
4. Run the benchmark
5. View the results
6. Destroy the instances

## Utilization

Once the CLI is installed you can use the benchmark tool by running the following command :

```bash
bench validate # validate the config
bench creds # check your cloud credentials
bench create # create the instances on the cloud
bench connection # check the connection to the instances
bench deploy # deploy the model on the instance
bench run # run the benchmark
bench results # view the results
bench destroy # destroy the instances
```

```mermaid
graph LR
    A[Start] --> B[Validate Config]
    B --> C[Check Cloud Credentials]
    C --> D[Create Instances]
    D --> E[Check Connection]
    E --> F[Deploy Model]
    F --> G[Run Benchmark]
    G --> H[View Results]
    H --> I[Destroy Instances]
    I --> J[End]
```

## Roadmap

- [ ] Run benchmarks on Scaleway
- [ ] Run benchmarks on GCP
- [ ] Use Ollama

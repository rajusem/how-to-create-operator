Sample python app that prints Current UTC time with Platform details.

### Run Main logic
You can run python code using following command locally.
```bash
python -m src.driver
```

### Build container image and push to registry

```
docker build -t <registry>/sample-python-app:latest ./  --platform linux/amd64 --no-cache
docker push  <registry>/sample-python-app:latest
```
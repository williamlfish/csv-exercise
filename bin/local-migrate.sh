  flyway migrate \
    -user=root \
    -password=password \
    -url=jdbc:postgresql://localhost:5432/processed_files \
    -locations=filesystem:migrations
use rocket::response::status::BadRequest;
use rocket::serde::json::{json, Json};
use rocket::config::SecretKey;
use rocket_cors::{AllowedOrigins, CorsOptions};
use rdkafka::{
    ClientConfig, producer::{FutureProducer, FutureRecord},
};
use std::env;

#[derive(rocket::serde::Deserialize)]
struct Data {
    name: String,
    album: String,
    year: String,
    rank: String,
}

#[rocket::post("/data", data = "<data>")]
async fn receive_data(data: Json<Data>) -> Result<String, BadRequest<String>> {
    let received_data = data.into_inner();

    // Convertir los datos a JSON
    let json_data = json!({
        "name": received_data.name,
        "album": received_data.album,
        "year": received_data.year,
        "rank": received_data.rank
    }).to_string();

    if let Err(e) = send_to_kafka(&json_data.to_string()).await {
        return Err(BadRequest(format!("Failed to send data to Kafka: {}", e)));
    }

    let response = json!({
        "message": format!("Received data Rust: Name: {}, Album: {}, Year: {}, Rank: {}",
                           received_data.name, received_data.album, received_data.year, received_data.rank)
    });

    Ok(response.to_string())
}

async fn send_to_kafka(data: &str) -> Result<(), rdkafka::error::KafkaError> {
    // Configurar el productor de Kafka
    let kafka_broker = env::var("KAFKABROKER").expect("KAFKABROKER not set");
    let kafka_topic = env::var("TOPIC").expect("TOPIC not set");

    // let producer_config = ClientConfig::new()
    //     .set("bootstrap.servers", &kafka_broker)
    //     .set("message.timeout.ms", "5000");

    // let producer: FutureProducer = producer_config
    //     .create()
    //     .expect("Failed to create Kafka producer");

    let producer: FutureProducer = ClientConfig::new()
        .set("bootstrap.servers", &kafka_broker)
        .set("message.timeout.ms", "5000")
        .create()
        .expect("Failed to create Kafka producer");

    match producer
        .send(
            FutureRecord::to(&kafka_topic)
                .payload(data)
                .key(&[]), // Clave vacía
            None,
        )
        .await
    {
        Ok((_delivery_status, _message)) => Ok(()),
        Err((kafka_error, _message)) => Err(kafka_error),
    }
}

#[rocket::main]
async fn main() {
    let secret_key = SecretKey::generate(); // Genera una nueva clave secreta

    // Configuración de opciones CORS
    let cors = CorsOptions::default()
        .allowed_origins(AllowedOrigins::all())
        .to_cors()
        .expect("failed to create CORS fairing");

    let config = rocket::Config {
        address: "0.0.0.0".parse().unwrap(),
        port: 8080,
        secret_key: secret_key.unwrap(), // Desempaqueta la clave secreta generada
        ..rocket::Config::default()
    };

    // Montar la aplicación Rocket con el middleware CORS
    rocket::custom(config)
        .attach(cors)
        .mount("/", rocket::routes![receive_data])
        .launch()
        .await
        .unwrap();
}

// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "SparkSDK",
    platforms: [.iOS(.v17), .macOS(.v14)],
    products: [
        .library(name: "SparkSDK", targets: ["SparkSDK"]),
    ],
    dependencies: [],
    targets: [
        .target(
            name: "SparkSDK",
            dependencies: [],
            swiftSettings: [.enableExperimentalFeature("StrictConcurrency")]
        ),
    ]
)

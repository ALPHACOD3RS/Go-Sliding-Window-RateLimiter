# Go Sliding Window Rate Limiter

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Overview

Go Sliding Window Rate Limiter is an advanced rate limiter implementation based on the sliding window algorithm. It provides a reliable and efficient way to limit the rate of incoming requests to your application.

## Features

- Sliding window algorithm: The rate limiter uses the sliding window algorithm to track and enforce rate limits.
- Precise rate limiting: The rate limiter allows you to set precise limits on the number of requests allowed within a specific time window.
- Easy integration: The rate limiter is designed to be easy to integrate into your Go applications.
- Thread-safe: The rate limiter is thread-safe, ensuring that it can be used in concurrent environments without issues.
- Highly customizable: The rate limiter provides various configuration options to tailor it to your specific needs.

## Installation

To install the Go Sliding Window Rate Limiter, use the following command:

- First You have to setup Redis client so make sure to use ur Address and password then you have done setuping

- go mod tidy

- go run main.go

## for testing using curl

- youse this command
  - for i in {1..10}; do curl http://localhost:8000/rate-limit; done

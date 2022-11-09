# dns53

Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily, and privately within a Route 53 Private Hosted Zone (PHZ).

Easily collaborate with a colleague by exposing your EC2 within a team VPC. You could even hook up a locally running application to a local k3d cluster using an ExternalName service during development. Once your EC2 is exposed, control how it is accessed through your EC2 security groups.

Written in Go, dns53 is incredibly small and easy to install.

<div>
    <video controls>
        <source src="./docs/static/dns53.webm" type="video/webm">
        <source src="./docs/static/dns53.mp4" type="video/mp4">
    </video>
    <sub>recorded using <a href="https://github.com/charmbracelet/vhs" target="_blank">VHS</a> 💜</sub>
</div>

## Badges

[![Build status](https://img.shields.io/github/workflow/status/purpleclay/dns53/ci?style=flat-square&logo=go)](https://github.com/purpleclay/dns53/actions?workflow=ci)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/purpleclay/dns53?style=flat-square)](https://goreportcard.com/report/github.com/purpleclay/dns53)
[![Go Version](https://img.shields.io/github/go-mod/go-version/purpleclay/dns53.svg?style=flat-square)](go.mod)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=purpleclay_dns53&metric=coverage)](https://sonarcloud.io/summary/new_code?id=purpleclay_dns53)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=purpleclay_dns53&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=purpleclay_dns53)

## Documentation

Check out the latest [documentation](https://purpleclay.github.io/dns53/)

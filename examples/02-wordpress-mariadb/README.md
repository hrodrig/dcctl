# Example 02: WordPress + MariaDB

Classic CMS stack: WordPress (Apache) and MariaDB. Good example of a "real world" LAMP-style setup managed with dcctl.

## Stack

- **MariaDB** — database for WordPress (persistent volume)
- **WordPress** — official WordPress image (Apache), port 8080

## Run with dcctl

From the **dcctl repo root**:

```bash
dcctl --config-file=examples/02-wordpress-mariadb/dcctl.yml up
```

Then open http://localhost:8080 and complete the WordPress setup wizard.

## Useful commands

```bash
dcctl --config-file=examples/02-wordpress-mariadb/dcctl.yml status
dcctl --config-file=examples/02-wordpress-mariadb/dcctl.yml logs -f wordpress
dcctl --config-file=examples/02-wordpress-mariadb/dcctl.yml down
```

## Files

- `dcctl.yml` — dcctl config
- `default/wordpress.yml` — Compose: MariaDB + WordPress, volume for DB data

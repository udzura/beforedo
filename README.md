# beforedo

Run before you do development

[![Build Status](https://travis-ci.org/udzura/beforedo.svg)](https://travis-ci.org/udzura/beforedo)

## Install

```bash
$ go get github.com/udzura/beforedo
```

Or [download here and put in under `$PATH`](https://github.com/udzura/beforedo/releases/latest)

## How to use

Set up `Before.yaml`:

```yaml
- task: mysql.server start
  port: 3306
- task: memcached -p 11211 -m 64m -d
  port: 11211
- task: bundle install --path vendor/bundle
  success: bundle check
- task: cp config/database.yml.sample config/database.yml
  file: config/database.yml
- task: bundle exec rake db:migrate
  always: true
- task: bundle exec rails s
  front: true
```

Then:

```bash
$ beforedo
```

## Task details

### `port`

* Skip task when specified port is already listened

### `file`

* Skip task when the destination file exists(if the destination is directory, all tasks fails)

### `success`

* Skip task when the specified command successfully exits

### `always`

* Always run the task

### `front`

* Always run the task, which is blocking to terminal. e.g. `rails s`, `npm start` &c.
  * You can run multipul tasks in front by `Procfile`

## License

[MIT](./LICENSE).


## Contributing

Pull Request welcome!

`-c` option is now available to specify config file(which is generally useful to debug).

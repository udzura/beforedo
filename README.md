# beforedo

Run before you do development


## How to use

Set up `Before.yaml`:

```
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
$ before do
```

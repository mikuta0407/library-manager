# library-manager

蔵書・CDを管理できるAPIサーバ


## これはなに
- 同人誌を主とする蔵書や、CDの管理をするためのAPIサーバです。
- 一覧取得・個別取得・追加・更新・削除・検索ができます

まだできないこと:
- Webフロントの表示

## 取得方法

### シングルバイナリで動かす場合

1. https://github.com/mikuta0407/library-manager/releases から最新の各種アーキテクチャに対応したバイナリをダウンロード
2. chmod +x等で権限付与
3. お好きなところへ配置
4. ```
    ./library-manager-<アーキテクチャ> <各種引数>
    ```

### ソースから動かす場合

(go installへの対応はいつかやります)

```console
https://github.com/mikuta0407/library-manager.git
cd library-manager
go run main.go <各種引数>
```

## コマンド

### サーバ起動

実行ファイル名は適宜読み替えて下さい

- APIサーバー実行
  - ```
    ./library-manager server [-l hostaddress] [-p listenaddess] [-f dbfilepath]
    ```
    - `-l` (`--httphost`)
      - 受けるホストのIP/ホスト名
      - 指定しない場合のデフォルトは`0.0.0.0`
    - `-p` (`--httpport`)
      - listenするポート番号
      - 指定しない場合のデフォルトは`8080`
    - `-f` (`--dbfile`)
      - 使用するSQLite3のDBファイルパス
      - 指定しない場合のデフォルトは`./library.db`
  - 例 (ローカルのみ8888番で受け、library-01.dbを使用する例)
    ```
    ./library-manager -l 127.0.0.1 -p 8888 -f ./library-01.db
    ```
- DBファイル生成
  - ```
    ./library-mamanger initdb [-f filepath]
    ```
    - `-f` (`--filepath`)
      - 出力ファイル名(指定しない場合のデフォルトは`./library.db`)

## API仕様

- /api/list: GET
    - 一覧
    - `/{カテゴリ名}` (`all`にすると全件取得)
    - 応答例 /api/list/book
        ```
        $ curl -sS localhost:8080/api/list/book | jq .
        {
          "items": [
            {
              "id": 1,
              "title": "testtitle",
              "author": "testauthor",
              "code": "12345678",
              "purchase": "C101",
              "place": "home",
              "note": "testnote",
              "image": "No Image""
            },
            {
              "id": 2,
              "title": "testtitle2",
              "author": "testauthor2",
              "code": "testauthor2",
              "purchase": "",
              "place": "234567",
              "note": "home2",
              "image": "https://example.com/hoge.jpg"
            },
            ...
          ]
        }
        ```
- /api/detail: GET
    - 詳細
    - `/{カテゴリ名}`
        - /{id}
    - 使用例 /api/detail/book/1
        ```
        $ curl -sS localhost:8080/api/detail/1 | jq .
        {
          "id": 1,
          "title": "testtitle",
          "category": "book"
          "author": "testauthor",
          "code": "12345678",
          "purchase": "C101",
          "place": "home",
          "note": "testnote",
          "image": null
        }
        ```
    - 存在しないIDの場合
        ```
        $ curl -sS localhost:8080/api/detail/100 | jq .
        {
          "message": "Not Found",
          "detail": "No record"
        }
- /api/search: POST
    - 検索
    - title/author/code/purchase/place/noteでLIKE検索する
    - 使用例
        ```
        $ curl -sS -X POST -H "Content-Type: application/json" -d  '{"title":"test","author":"test"}' http://localhost:8080/api/search | jq .
        {
          "items": [
            {
              "id": 1,
              "title": "testtitle",
              "author": "testauthor",
              "code": "12345678",
              "purchase": "C101",
              "place": "home",
              "note": "testnote",
              "image": null
            },
            {
              "id": 2,
              "title": "testtitle2",
              "author": "testauthor2",
              "code": "testauthor2",
              "purchase": "",
              "place": "234567",
              "note": "home2",
              "image": "Ng=="
            },
            ...
          ]
        }
        ```
    - 存在しない場合
        ```
        $ curl -sS -X POST -H "Content-Type: application/json" -d  '{"title":"あああ","author":"あああ"}' http://localhost:8080/api/search | jq .
        {
          "items": null
        }
- /api/create: POST
    - レコード作成
    - RequestBody: `{"title":"hoge","artist":"fuga".....}`
        - Response: id
    - title以外は空欄でもOK
    - **注意**
        - 登録時にAPIサーバーはコード重複を検知しない
        - もし重複を避けたい場合は、クライアントは一度searchエンドポイントで検索を行い、存在有無を確認する必要がある。
            - searchしてからcreateする
    - 使用例
        ```
        $ curl -X POST -H "Content-Type: application/json" \
        -d '{"title":"テストのタイトル-あいうえお","category":"book","author":"テストの著者-abcd","code":"abcd1234","purchase":"C101","place":"倉庫","note":" にゃーん"}' \
        http://localhost:8080/api/create
        {"message":"Success","id":"10"}
        $ curl -sS localhost:8080/api/detail/10 | jq .
        {
          "id": 10,
          "title": "テストのタイトル-あいうえお",
          "author": "テストの著者-abcd",
          "code": "abcd1234",
          "purchase": "C101",
          "place": "倉庫",
          "note": "にゃーん",
          "image": null
        }
        ```
- /api/update: PUT
    - レコード編集
    - /{id}
    - RequestBody: `{"title":"hoge","artist":"fuga".....}`
        - Response: id
    - **注意**
        - RequestBodyの内容にすべて上書きするので、更新したい場合でも全項目を付与する必要がある
    - 使用例
        ```
        curl -X PUT -H "Content-Type: application/json" \
        -d '{"title":"テストのタイトル-あいうえお","author":"テストの著者-abcd","code":"abcd1234","purchase":"C101","place":"自宅","note":"20XX/MM/DD 自宅へ移動"}' \
        http://localhost:8080/update/10
        {"message":"Success","id":"10"}
        $ curl -sS localhost:8080/api/detail/10 | jq .
        {
          "id": 10,
          "title": "テストのタイトル-あいうえお",
          "author": "テストの著者-abcd",
          "code": "abcd1234",
          "purchase": "C101",
          "place": "自宅",
          "note": "20XX/MM/DD 自宅へ移動",
          "image": null
        }
        ```
    - 存在しないIDへUPDATE要求した場合は404が帰ってくる
        ```
        $ curl -X PUT -H "Content-Type: application/json" \
        -d '{"title":"テストのタイトル-あいうえお"}' \
        http://localhost:8080/update/100
        {"message":"Not Found","detail":"No record"}
        ```
- /api/delete: DELETE
    - レコード削除
    - /{id}
    - 使用例
        ```
        $ curl -sS -X DELETE http://localhost:8080/api/delete/10 | jq .
        {
          "message": "Success",
          "id": "10"
        }
        $ curl -sS localhost:8080/detail/10 | jq . 
        {
          "message": "Not Found",
          "detail": "No record"
        }
        ```

### メソッドが異なった場合
以下のような応答となり、処理が行われない
```
$ curl -sS http://localhost:8080/delete/10 | jq .
{
  "message": "Method Not Allowed",
  "detail": "Use DELETE Method"
}
```


## DB設計

### library
|カラム名|型|内容|
|:-:|:-:|:-:|
|id|INTEGER|AUTOINCREMENTするやつ|
|title|TEXT|題名|
|author|TEXT|著者|
|code|TEXT|バーコードのやつ(CODE39の場合アルファベットがある)|
|purchase|TEXT|購入場所|
|place|TEXT|存在場所|
|note|TEXT|備考欄|
|image|TEXT|URL|

```sql
CREATE TABLE "library" (
	"id"	INTEGER,
	"title"	TEXT NOT NULL,
	"category" TEXT,
	"author"	TEXT,
	"code"	TEXT,
	"purchase" TEXT,
	"place"	TEXT,
	"note"	TEXT,
	"image"	TEXT,
	"user" TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
)
```

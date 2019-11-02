# gortfolio

<h2>概要</h2>
<p>golangで作ったポートフォリオ（portfolio） 略して、gortfolio</p>
<p>チャット、スクレイピング、しりとりなど、種々雑多な機能を詰め込んだサイト</p>
<h2>機能一覧</h2>
<ol>
    <li>ページ別アクセス数の棒グラフ表示</li>
    <li>Qiitaのトレンド記事表示（スクレイピング,goquery）</li>
    <li>認証機能</li>
    <ul>
      <li>Googleアカウントによるログイン(gomniauth)</li>
      <li>新規登録、ログイン、ログアウト</li>
      <li>テストユーザーでログイン</li>
    </ul>
    <li>websocketを用いたチャット機能</li>
    <ul>
      <li>メッセージ送信</li>
      <li>画像のアップロード</li>
    </ul>
    <li>一人しりとり機能</li>
    <li>タスク管理機能</li>
    <ul>
      <li>タスクの作成・編集・削除</li>
    </ul>
    <li>あしあと機能</li>
    <li>QRコード表示</li>
</ol>  

<h2>技術</h2>
<ol>
    <li>Golang1.13, gorm, Go Modules</li>
    <li>GCE</li>
    <li>SQLite3</li>
    <li>Bootstrap4</li>
</ol>
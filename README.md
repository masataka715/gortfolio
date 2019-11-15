<h1>gortfolio</h1>

<h2>概要</h2>
<p>golangで作ったポートフォリオ（portfolio） 略して、gortfolio</p>
<p>ブラックジャック、条文検索、チャット、スクレイピングなど、種々雑多な機能を詰め込んだサイトを、GCEへデプロイ</p>

<h2>機能</h2>
<ol>
    <li>現在の東京の天気・気温・風速表示（OpenWeather API)</li>
    <li>ブラックジャック機能</li>
    <ul>
        <li>ディーラーと対戦し、勝敗を記録できます</li>
    </ul>
    <li>刑法の条文検索機能</li>
    <ul>
        <li>全文の表示（スクレイピング,goquery）</li>
        <li>条文番号（算用数字）による条文検索</li>
    </ul>
    <li>あしあと機能</li>
    <ul>
        <li>ページ別アクセス数の棒グラフ表示（gonum）</li>
        <li>閲覧者のアクセスページ・日時の表示・PDF出力機能（gopdf）</li>
    </ul>
    <li>Qiitaのトレンド記事表示（スクレイピング,goquery）</li>
    <li>認証機能</li>
    <ul>
      <li>Googleアカウントによるログイン(gomniauth)</li>
      <li>新規登録、ログイン、ログアウト</li>
      <li>テストユーザーでログイン</li>
      <li>フラッシュメッセージ表示</li>
    </ul>
    <li>websocketを用いたチャット機能</li>
    <ul>
      <li>メッセージ送信</li>
      <li>画像のアップロード、名前変更</li>
    </ul>
    <li>一人しりとり機能</li>
    <li>タスク管理機能</li>
    <ul>
      <li>タスクの作成・編集・削除</li>
    </ul>
    <li>QRコード表示</li>
</ol>  

<h2>技術</h2>
<ol>
    <li>Golang1.13, gorm, Go Modules</li>
    <li>GCE, Cloud DNS</li>
    <li>SQLite3</li>
    <li>Bootstrap4</li>
</ol>


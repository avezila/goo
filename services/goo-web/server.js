var webpack = require('webpack');
var WebpackDevServer = require('webpack-dev-server');
var config = require('./webpack.config');

console.log(process.env.SERVER)

global.ENV = {
  SERVER : process.env.SERVER
}

new WebpackDevServer(webpack(config), {
  publicPath: config.output.publicPath,
  hot: true,
  historyApiFallback: true,
}).listen(process.env.PORT, function (err, result) {
  if (err) {
    console.log(err);
  }

  console.log('Listening at localhost:'+process.env.PORT);
});

var gulp = require('gulp');
var coffee = require('gulp-coffee');
var coffeelint = require('gulp-coffeelint');
var concat = require('gulp-concat');
var compress = require('gulp-yuicompressor');
var sourcemaps = require('gulp-sourcemaps');
var jasmine = require('gulp-jasmine-phantom');
var plumber = require('gulp-plumber');
var del = require('del');
var order = require('gulp-order');
var process = require('child_process');

var paths = {
  tmp: 'bin/tmp/',
  dev: 'bin/dev/',
  release: 'bin/release/',
  spec: 'spec/',
  assets: {
    coffee: 'src/scripts/*.coffee',
    widgets: 'src/scripts/widgets/*.coffee',
    models: '',
    html: 'src/*.html',
    templates: 'src/templates/*.tmpl',
    css: 'src/styles/*.css',
    vendor: {
      js: 'src/scripts/vendor/*.js',
      css: 'src/styles/vendor/*.css'
    }
  }
};

var STATE_OK = 0,
    STATE_ERR = 1;

var state = STATE_OK;

function handleError (err) {
  console.log(err.toString());
  this.emit('end');
  state = STATE_ERR;
  setTimeout(function() {
    state = STATE_OK;
  }, 200)
}

gulp.task('clean', function() {
  return del(['bin']);
});

gulp.task('compile:assets:html', ['clean'], function() {
  if (state === STATE_OK) {
      gulp.src(paths.assets.html)
      .pipe(plumber(handleError))
      .pipe(concat('index.html'))
      .pipe(gulp.dest(paths.dev));

      return gulp.src(paths.assets.html)
        .pipe(plumber(handleError))
        .pipe(concat('index.html'))
        .pipe(gulp.dest(paths.release));
  } else if (state === STATE_ERR) {
    console.log('State is "error" skipping...');
  }
});

gulp.task('compile:assets:templates', ['compile:assets:html'], function() {
  if (state === STATE_OK) {
      gulp.src(paths.assets.templates)
      .pipe(plumber(handleError))
      .pipe(gulp.dest(paths.dev));

      return gulp.src(paths.assets.templates)
        .pipe(plumber(handleError))
        .pipe(gulp.dest(paths.release));
  } else if (state === STATE_ERR) {
    console.log('State is "error" skipping...');
  }
});

gulp.task('compile:assets:css', ['compile:assets:templates'], function() {
  if (state === STATE_OK) {
      gulp.src([paths.assets.vendor.css, paths.assets.css])
      .pipe(order([paths.assets.vendor.css, paths.assets.css], {base: './'}))
      .pipe(plumber(handleError))
      .pipe(concat('ed.css'))
      .pipe(gulp.dest(paths.dev));

      return gulp.src([paths.assets.vendor.css, paths.assets.css])
        .pipe(order([paths.assets.vendor.css, paths.assets.css], {base: './'}))
        .pipe(plumber(handleError))
        .pipe(concat('ed.css'))
        .pipe(compress({
          type: 'css'
        }))
        .pipe(gulp.dest(paths.release));
  } else if (state === STATE_ERR) {
    console.log('State is "error" skipping...');
  }
});

gulp.task('compile:coffee', ['compile:assets:css'], function() {
  if (state === STATE_OK) {
    return gulp.src([paths.assets.coffee, paths.assets.widgets])
      .pipe(order([paths.assets.widgets, paths.assets.coffee], {base: './'}))
    	.pipe(plumber(handleError))
    	//.pipe(sourcemaps.init())
    	.pipe(coffeelint())
      .pipe(coffeelint.reporter())
    	.pipe(coffee())
      .pipe(concat('ed.coffee.js'))
      //.pipe(sourcemaps.write())
      .pipe(gulp.dest(paths.tmp));
  } else if (state === STATE_ERR) {
    console.log('State is "error" skipping...');
  }
});

gulp.task('compile:dev', ['compile:coffee'], function() {
  if (state === STATE_OK) {
    return gulp.src([paths.assets.vendor.js, paths.tmp + '*.coffee.js'])
      .pipe(order([paths.assets.vendor.js, paths.tmp + '*.coffee.js'], {base: './'}))
      .pipe(plumber(handleError))
      .pipe(concat('ed.js'))
      .pipe(gulp.dest(paths.dev));
  } else if (state === STATE_ERR) {
    console.log('State is "error" skipping...');
  }
});

gulp.task('compile:release', ['compile:dev'], function() {
  if (state === STATE_OK) {
    return gulp.src(paths.dev + '*.js')
    	.pipe(plumber(handleError))
      .pipe(compress({
        type: 'js'
      }))
      .pipe(concat('ed.js'))
      .pipe(gulp.dest(paths.release));
  } else if (state === STATE_ERR) {
    console.log('State is "error" skipping...');
  }
});

// gulp.task('compile:spec', ['compile:release'], function() {
//   if (state === STATE_OK) {
//   	return gulp.src(paths.spec)
//   		.pipe(plumber(handleError))
//   		.pipe(coffeelint())
//       .pipe(coffeelint.reporter())
//     	.pipe(coffee())
//     	.pipe(concat('malefic.spec.js'))
//     	.pipe(gulp.dest('bin/spec'));
//   } else if (state === STATE_ERR) {
//     console.log('State is "error" skipping...');
//   }
// });

// gulp.task('spec', ['compile:spec'], function() {
//   if (state === STATE_OK) {
//     return gulp.src(['bin/dev/malefic.js', 'bin/spec/malefic.spec.js'])
//     	.pipe(jasmine({
//     		integration: true
//     	}));
//   } else if (state === STATE_ERR) {
//     console.log('State is "error" skipping...');
//   }
// });

gulp.task('docker', function(event) {
  var docker = process.spawn('./build.sh');
  docker.stdout.on('data', function(data) {
   console.log(data.toString());
  });
  docker.stdout.on('end', function() {
  });

  docker.on('exit',function(code) {
    console.log('./build.sh - ' + code);
  });
})

gulp.task('watch', function(event) {
  gulp.watch([paths.assets.coffee, paths.assets.widgets, paths.assets.templates, paths.assets.html, paths.assets.css], ['docker']);
});

gulp.task('default', ['compile:release']);

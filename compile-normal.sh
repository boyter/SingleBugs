cat ./assets/tmpl/base.html ./assets/tmpl/jquery.underscore.html ./assets/tmpl/javascript.start.html ./assets/tmpl/custom.js ./assets/tmpl/javascript.end.html > ./tmpl/base.html && cat ./assets/tmpl/styles.css > ./tmpl/styles.css && go build fizzbuzz.go && ./fizzbuzz
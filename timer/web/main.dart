import 'dart:html';
import 'dart:async';

void main() {
	var m = Uri.base.queryParameters['m'];
	var t = Uri.base.queryParameters['t'];

	if (m == "" || m == null) {
		m = "開啟網頁";
	}

	DateTime timestamp;
	if (t == "" || t == null) {
		timestamp = new DateTime.now();
	} else {
		timestamp = DateTime.tryParse(t);
		if (timestamp == null) {
			timestamp = new DateTime.now();
		}
	}

	querySelector('#title').text = m;
	Timer.periodic(Duration(milliseconds: 200), (Timer timer) {
		var now = new DateTime.now();
		var diff = now.difference(timestamp);

		if (diff.isNegative) {
			diff *= -1;
			querySelector('#tense').text = '還有';
		} else {
			querySelector('#tense').text = '過了';
		}
		var msg = '${diff.inDays} 天 ${diff.inHours%24} 時 ${diff.inMinutes%60} 分 ${diff.inSeconds%60} 秒';
		querySelector('#timer').text = msg;
	});
}

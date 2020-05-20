import 'dart:html';
import 'dart:async';

void NewTitleEvent(KeyboardEvent e) {
	if (e.charCode == KeyCode.ENTER) {
		var title = querySelector('#title').text;

		var args = Uri.parse(window.location.href).queryParameters;
		window.location = '/?m=$title&t=${args["t"] ?? ""}';
	}
}

void ButtenEvent(Event e) {
	if (querySelector('#btn').text == '🔓') {
		querySelector('#btn').text = '🔒';
		querySelector('#timestamp').classes.add('hidden');

		var dom = (querySelector('#timestamp') as InputElement);
		if (dom.dataset['src'] != null && dom.value != dom.dataset['src']) {
		var args = Uri.parse(window.location.href).queryParameters;
		window.location = '/?m=${args["m"] ?? ""}&t=${dom.value}';
		}
	} else {
		querySelector('#btn').text = '🔓';
		querySelector('#timestamp').classes.remove('hidden');
	}
}

void main() {
	var timestamp = DateTime.tryParse(querySelector('#timer').getAttribute('data-src'));
	if (timestamp == null) {
		timestamp = new DateTime.now();
	}

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

	/* ---- event ---- */
	querySelector('#title').onKeyPress.listen(NewTitleEvent);
	querySelector('#btn').onClick.listen(ButtenEvent);
}

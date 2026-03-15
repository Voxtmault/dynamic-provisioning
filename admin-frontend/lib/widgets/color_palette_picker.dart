import 'package:flutter/material.dart';
import 'package:flutter_colorpicker/flutter_colorpicker.dart';

/// Three colour slots: primary (index 0), secondary (index 1), background (index 2).
class ColorPalettePicker extends StatefulWidget {
  final List<String> initialPalette; // list of hex strings, may be empty
  final ValueChanged<List<String>> onChanged;

  const ColorPalettePicker({
    super.key,
    required this.initialPalette,
    required this.onChanged,
  });

  @override
  State<ColorPalettePicker> createState() => _ColorPalettePickerState();
}

class _ColorPalettePickerState extends State<ColorPalettePicker> {
  static const _labels = ['Primary', 'Secondary', 'Background'];
  static const _defaults = [
    Color(0xFF6200EE),
    Color(0xFF03DAC6),
    Color(0xFFF5F5F5),
  ];

  late final List<Color> _colors;

  @override
  void initState() {
    super.initState();
    _colors = List.generate(3, (i) {
      final hex = widget.initialPalette.length > i
          ? widget.initialPalette[i]
          : null;
      return _fromHex(hex) ?? _defaults[i];
    });
  }

  void _notify() {
    widget.onChanged(_colors.map(_toHex).toList());
  }

  Future<void> _pick(int index) async {
    Color picked = _colors[index];

    await showDialog<void>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('${_labels[index]} colour'),
        content: SingleChildScrollView(
          child: HueRingPicker(
            pickerColor: picked,
            onColorChanged: (c) => picked = c,
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.of(ctx).pop(),
            child: const Text('Select'),
          ),
        ],
      ),
    );

    setState(() => _colors[index] = picked);
    _notify();
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      children: List.generate(3, (i) {
        return Expanded(
          child: Padding(
            padding: EdgeInsets.only(right: i < 2 ? 8 : 0),
            child: _SlotButton(
              label: _labels[i],
              color: _colors[i],
              hexString: _toHex(_colors[i]),
              onTap: () => _pick(i),
            ),
          ),
        );
      }),
    );
  }

  static String _toHex(Color c) {
    final r = c.r.round();
    final g = c.g.round();
    final b = c.b.round();
    return '#${r.toRadixString(16).padLeft(2, '0')}${g.toRadixString(16).padLeft(2, '0')}${b.toRadixString(16).padLeft(2, '0')}'.toUpperCase();
  }

  static Color? _fromHex(String? hex) {
    if (hex == null || hex.isEmpty) return null;
    final cleaned = hex.replaceFirst('#', '');
    final value = int.tryParse(
      cleaned.length == 6 ? 'FF$cleaned' : cleaned,
      radix: 16,
    );
    return value != null ? Color(value) : null;
  }
}

class _SlotButton extends StatelessWidget {
  final String label;
  final Color color;
  final String hexString;
  final VoidCallback onTap;

  const _SlotButton({
    required this.label,
    required this.color,
    required this.hexString,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Column(
        children: [
          Container(
            height: 48,
            decoration: BoxDecoration(
              color: color,
              borderRadius: BorderRadius.circular(8),
              border: Border.all(color: Colors.black12),
            ),
          ),
          const SizedBox(height: 4),
          Text(
            label,
            style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w600),
          ),
          Text(
            hexString,
            style: const TextStyle(fontSize: 10, color: Colors.black45),
          ),
        ],
      ),
    );
  }
}

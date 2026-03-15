import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:provider/provider.dart';

import '../providers/auth_provider.dart';
import '../providers/tenant_provider.dart';
import '../widgets/color_palette_picker.dart';
import 'success_screen.dart';

class RegisterTenantScreen extends StatefulWidget {
  const RegisterTenantScreen({super.key});

  @override
  State<RegisterTenantScreen> createState() => _RegisterTenantScreenState();
}

class _RegisterTenantScreenState extends State<RegisterTenantScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();

  List<String> _palette = [];
  Uint8List? _photoBytes;
  String? _photoMime;
  String? _photoName;
  bool _isLoading = false;
  String? _errorMessage;

  @override
  void dispose() {
    _nameCtrl.dispose();
    super.dispose();
  }

  Future<void> _pickPhoto() async {
    final picker = ImagePicker();
    final file = await picker.pickImage(
      source: ImageSource.gallery,
      maxWidth: 1024,
      maxHeight: 1024,
      imageQuality: 85,
    );
    if (file == null) return;

    final bytes = await file.readAsBytes();
    setState(() {
      _photoBytes = bytes;
      _photoName = file.name;
      _photoMime = file.mimeType ?? 'image/jpeg';
    });
  }

  Future<void> _submit() async {
    if (!(_formKey.currentState?.validate() ?? false)) return;
    if (_palette.length < 3) {
      setState(() => _errorMessage =
          'Please select all three colours (Primary, Secondary, Background).');
      return;
    }

    setState(() {
      _isLoading = true;
      _errorMessage = null;
    });

    try {
      final token = context.read<AuthProvider>().token!;
      final resp = await context.read<TenantProvider>().registerTenant(
            token: token,
            name: _nameCtrl.text.trim(),
            colorPalette: _palette,
            photoBytes: _photoBytes,
            photoMimeType: _photoMime,
          );

      if (!mounted) return;
      Navigator.of(context).pushReplacement(
        MaterialPageRoute<void>(
          builder: (_) => SuccessScreen(response: resp),
        ),
      );
    } catch (e) {
      final msg = e.toString();
      setState(() {
        _isLoading = false;
        _errorMessage =
            msg.startsWith('Exception: ') ? msg.substring(11) : msg;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Register New Tenant')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Center(
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 560),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  // ── Name ──────────────────────────────────────────────────
                  TextFormField(
                    controller: _nameCtrl,
                    textInputAction: TextInputAction.next,
                    decoration: const InputDecoration(
                      labelText: 'Tenant Name',
                      hintText: 'e.g. Acme Corp',
                      prefixIcon: Icon(Icons.business),
                      border: OutlineInputBorder(),
                    ),
                    validator: (v) {
                      if (v == null || v.trim().isEmpty) {
                        return 'Tenant name is required';
                      }
                      if (v.trim().length > 255) {
                        return 'Name must be 255 characters or fewer';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 28),

                  // ── Colour palette ─────────────────────────────────────────
                  const Text(
                    'Colour Palette',
                    style:
                        TextStyle(fontWeight: FontWeight.w600, fontSize: 14),
                  ),
                  const SizedBox(height: 4),
                  const Text(
                    'Tap a swatch to choose a colour.',
                    style: TextStyle(fontSize: 12, color: Colors.black45),
                  ),
                  const SizedBox(height: 10),
                  ColorPalettePicker(
                    initialPalette: _palette,
                    onChanged: (p) => setState(() => _palette = p),
                  ),
                  const SizedBox(height: 28),

                  // ── Photo ─────────────────────────────────────────────────
                  const Text(
                    'Profile Photo (optional)',
                    style:
                        TextStyle(fontWeight: FontWeight.w600, fontSize: 14),
                  ),
                  const SizedBox(height: 10),
                  _PhotoPicker(
                    photoBytes: _photoBytes,
                    photoName: _photoName,
                    onPick: _pickPhoto,
                    onClear: () =>
                        setState(() {
                          _photoBytes = null;
                          _photoName = null;
                          _photoMime = null;
                        }),
                  ),

                  // ── Error ─────────────────────────────────────────────────
                  if (_errorMessage != null) ...[
                    const SizedBox(height: 20),
                    Container(
                      padding: const EdgeInsets.all(10),
                      decoration: BoxDecoration(
                        color: Colors.red.shade50,
                        borderRadius: BorderRadius.circular(6),
                        border: Border.all(color: Colors.red.shade200),
                      ),
                      child: Text(
                        _errorMessage!,
                        style: TextStyle(
                            color: Colors.red.shade800, fontSize: 13),
                      ),
                    ),
                  ],

                  const SizedBox(height: 28),

                  // ── Submit ────────────────────────────────────────────────
                  FilledButton(
                    onPressed: _isLoading ? null : _submit,
                    child: _isLoading
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              color: Colors.white,
                            ),
                          )
                        : const Text('Register Tenant'),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

// ── Photo picker sub-widget ────────────────────────────────────────────────────

class _PhotoPicker extends StatelessWidget {
  final Uint8List? photoBytes;
  final String? photoName;
  final VoidCallback onPick;
  final VoidCallback onClear;

  const _PhotoPicker({
    required this.photoBytes,
    required this.photoName,
    required this.onPick,
    required this.onClear,
  });

  @override
  Widget build(BuildContext context) {
    if (photoBytes != null) {
      return Row(
        children: [
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: Image.memory(
              photoBytes!,
              width: 72,
              height: 72,
              fit: BoxFit.cover,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              photoName ?? 'Selected photo',
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontSize: 13),
            ),
          ),
          IconButton(
            icon: const Icon(Icons.close, color: Colors.red),
            tooltip: 'Remove photo',
            onPressed: onClear,
          ),
        ],
      );
    }

    return OutlinedButton.icon(
      onPressed: onPick,
      icon: const Icon(Icons.photo_library_outlined),
      label: const Text('Choose from Gallery'),
    );
  }
}

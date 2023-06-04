
<div dir="rtl">
سلام.

تو این ریپو قصد داریم در ابتدا یک کلاستر ربیت روی کوبر بالا بیاریم. ربیت‌ام‌کیو یک نوعی از سرویسه که بهمون این امکان رو می‌ده که پیام‌ها رو داخل تعدادی صف قرار بدیم. در انتهای کار پس از بالا آوردن ربیت ام‌کیو در مورد استفادش بیشتر صحبت می‌کنیم.
 برای انجام این کار ما یمل‌های مربوط به کوبر رو خودمون می‌نویسیم و هدف بیشتر آشنایی با کامپوننت‌های کوبر و نحوه‌ی کارکرد ربیت‌ام‌کیو کنار هم هستش. توی استفاده‌ی پروداکشنی بهتره که از هلم چارت‌های مخصوص به ربیت‌ام‌کیو استفاده کنید.

برای شروع کار ما فرض می‌کنیم یک کلاستر کوبر به همراه کلاینت کوبر رو در اختیار دارید.

در ابتدا ما یک کانفیگ مپ می‌سازیم و فایل‌های کانفیگ مربوط به ربیت‌ام‌کیو رو توی اون قرار می‌دیم.

</div>

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
 name: rabbitmq-config
 namespace: rabbitmq-golang
data:
 enabled_plugins: |
  [rabbitmq_federation,rabbitmq_management,rabbitmq_peer_discovery_k8s].
 rabbitmq.conf: |
  log.console = true
  cluster_formation.peer_discovery_backend  = rabbit_peer_discovery_k8s
  cluster_formation.k8s.host = kubernetes.default.svc.cluster.local
  cluster_formation.k8s.address_type = hostname
  cluster_formation.k8s.service_name = rabbitmq-headless
  cluster_formation.node_cleanup.interval = 90
  cluster_formation.node_cleanup.only_log_warning = true
  cluster_formation.discovery_retry_limit = 60
  cluster_formation.discovery_retry_interval = 2000
  load_definitions = /etc/rabbitmq/import.json
 init-config.sh: |
  cp /tmp/config/rabbitmq.conf /etc/rabbitmq/rabbitmq.conf
  cp /tmp/config/enabled_plugins /etc/rabbitmq/enabled_plugins
  echo $RABBITMQ_ERLANG_COOKIE > /var/lib/rabbitmq/.erlang.cookie
  chmod 600 /var/lib/rabbitmq/.erlang.cookie
  envsubst < /tmp/config/import.json > /etc/rabbitmq/import.json
 import.json: |
  {
    "permissions": [
      {
        "configure": ".*",
        "read": ".*",
        "user": "admin",
        "vhost": "test",
        "write": ".*"
      }
    ],
    "policies": [
      {
        "apply-to": "queues",
        "definition": {
          "ha-mode": "exactly",
          "ha-params": 2,
          "ha-sync-mode": "automatic"
        },
        "name": "ha-fed-test",
        "pattern": ".*",
        "priority": 0,
        "vhost": "test"
      }
    ],
    "users": [
      {
        "hashing_algorithm": "rabbit_password_hashing_sha256",
        "limits": {},
        "name": "admin",
        "password": "$RABBITMQ_ADMIN_PASSWORD",
        "tags": [
          "administrator"
        ]
      }
    ],
    "vhosts": [
      {
        "limits": [],
        "metadata": {
          "description": "test vhost",
          "tags": []
        },
        "name": "test"
      }
    ]
  }
```

<div dir="rtl">
خب حالا این فایلا چی هستند. فایل اول یه فایله که توش پلاگین‌هایی که فعال هستند رو مشخص می‌کنیم. ما فعلا سه تا پلاگین فعلا کردیم. پلاگین <code>rabbitmq_federation</code> برای فعال کردن رپلیکیشن پیام‌های بین چند تا نود هست. هدفمون اینه که الان دو تا نود ربیت بیاریم بالا و پیام‌هامون رو به هر کدوم که فرستادیم خودش تو اون یکی هم بنویسه و بتونیم از هر دو تا نود پیام‌ها رو بگیریم. می‌تونیم در ادامه برای پخش کردن لود از چندین تا صف هم استفاده کنیم اگر حجم پیام‌هامون خیلی زیاد شد. فعلا دو تا نود میاریم که بتونیم اطمینان داشته باشیم که اگه یکیشون هم رفت پایین هنوز سرویس ربیتمون در دسترس باشه. پلاگین بعدی برای استفاده از رول ادمین تو ui ربیت هست و نکته‌ی خاصی نداره.
پلاگین آخر <code>rabbitmq_peer_discovery_k8s</code> هست. این پلاگین بهمون این قابلیت رو می‌ده که هر نسخه از ربیت که بالا میاد چطوری بقیه نسخه‌های خودش رو پیدا کنه و به هم دیگه جوین بشند. تو این پلاگین ربیت از api server کوبر می‌پرسه و لیست پاد‌های مربوط به خودش رو می‌گیره و اینجوری بقیه رو پیدا می‌کنه.

فایل بعدی `rabbit.conf`  هست. این فایل کانفیگ خود ربیت هستش و می‌تونید هر مدل کانفیگی که نیاز بود رو اضافه کنید. می‌تونید از داک خود ربیت‌ام‌کیو در [این آدرس](https://www.rabbitmq.com/configure.html) برای این کار استفاده کنید.

 ما در اینجا فقط کانفیگ‌های مربوط به پلاگین دیسکاوری کوبر رو گذاشتیم. اول اینکه حواستون به آدرس api سرور باشه که ممکنه عوض شده باشه در بعضی از کوبرها. بعد اون ما در کانفیگ‌ها قرار دادیم که این پلاگین از هاست‌نیم‌های پاد‌های دیگه استفاده کنه و برای این کار نیازه که سرویس لیست اندپوینت‌های rabbitmq-headless رو از api سرور بپرسه و اینجوری می‌تونه بقیه پاد ها رو پیدا کنه. این سرویس رو جلوتر می‌سازیم.

بعد اون می‌رسیم به فایل `import.json`. این فایل به این درد می‌خوره که کانفیگ‌های مربوط به ربیت از جمله بوزر‌ها و اطلاعات صف‌ها رو به عنوان کانفیگ داشته باشیم و نیازی نباشه اینا رو به صورت دستی از توی ui بسازیم. تو این فایل جیسون کار خاصی انجام ندادیم. به صورت کلی یک یوزر تست ساختیم با یک وی‌هاست تستی و به اون یوزر پرمیشن مربوط به این وی‌هاست رو دادیم. برای وی‌هاست تستیمون هم یه پالیسی گذاشتیم که پیام‌ها رو دو بار رپلیکیت کنه.

در آخر می‌رسیم به اسکریپت `init-config.sh`. این اسکریپت قراره قبل از بالا اومدن ربیت توی پاد اجرا بشه و فایل‌های کانفیگ رو آماده کنه. این کار برای این هست که ربیت نیاز داره که پرمیشن‌ فایل‌های کانفیگ ریدآنلی نباشه ولی وقتی این فایل‌ها رو توی پاد مانت کنیم ریدآنلی می‌شن و برای همین با این اسکریپت این فایلا رو کپی می‌کنیم یه جای دیگه. بعد از اون یه فایلی باید برای `erlang.cookie` بسازیم. این یک مقدار هست که توی کلاستر مقدارش ثابته و پاد‌های ربیت از این طریق می‌فهمن که هر دوشون توی یه کلاسترن و به همدیگه اعتماد می‌کنن. یه جور کلید مشترکه و برای مکانیم‌های امنیتی هستش. ما این مقدار رو وریبل می‌کنیم و از طریق این اسکریپت تو فایلش قرار می‌دیم و بهش پرمیشن مناسب می‌دیم. در نهایت فایل `import.json` رو آماده می‌کنیم. از اونجا پسوردها رو مستقیم نمی‌خواییم توی این فایل بنویسیم پس به صورت یه وریبل می‌ذاریم که در هنگام اجرا مثل یه تمپلیت پرش کنیم.

در قدم بعدی سیکرتی که شامل پسورد ادمین و ارلنگ کوکی هست رو می‌سازیم.
</div>


```yaml
apiVersion: v1
kind: Secret
metadata:
  name: rabbitmq-secret
  namespace: rabbitmq-golang
type: Opaque
data:
  RABBITMQ_ERLANG_COOKIE: 
  RABBITMQ_ADMIN_PASSWORD: 
```

<div dir="rtl">

این مقادیر رو برای تست به همین شکل استفاده می‌کنیم. در قدم بعدی به فایل  `rabc.yml` می‌رسیم. تو این فایل یه سرویس اکانت و یه رول و یه رول بایندیگ می‌سازیم. این کارا برای این هست که به پاد‌های ربیت این اجازه رو بدیم که بتونن از api سرور لیست هاست پاد‌های دیگه رو بگیرن و همه دیگه رو پیدا کنند.  

</div>

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: rabbitmq-service-account
  namespace: rabbitmq-golang

---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: rabbitmq-role
  namespace: rabbitmq-golang
rules:
  - apiGroups:
      - ""
    resources:
      - endpoints
    verbs:
      - get
      - list
      - watch
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: rabbitmq-role-binding
  namespace: rabbitmq-golang
subjects:
  - kind: ServiceAccount
    name: rabbitmq-service-account
    namespace: rabbitmq-golang
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: rabbitmq-role
```

<div dir="rtl">

همونطور که قبلتر هم گفتیم اینجا یک سرویس اکانت ساختیم و به این سرویس اکانت این دسترسی رو دادیم که اندپوینت‌ها رو لیست کنه و بگیره.

در ادامه می‌خواییم یه استیت فول ست بسازیم. استیت‌فول ست این امکان رو بهمون می‌ده که برعکس دیپلویمنت که یه سرویس بدون استیت هست، یه سرویس بسازیم که استیت رو داخل خودش نگه داره. ساده تر بخواییم نگاه کنیم این قابلیت رو بهمون می‌ده که به ازای هر پاد از رپلیکاست‌مون یه pvc هم بهش اتچ کنیم. pvc یک کامپوننت کوبرنتیزه که این امکان رو بهمون یه سری دیتا رو توش پرسیست کنیم و پایین رفتن پاد از بین نره. یه چیزی شبیه والیوم داکر ولی توی کلاستر کوبرنتیز این قابلیت رو داره که به نود‌های مجزا وصل بشه. یه همچین یملی رو می‌نویسیم.
</div>

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
 name: rabbitmq
 namespace: rabbitmq-golang
spec:
 serviceName: rabbitmq-headless
 replicas: 2
 selector:
  matchLabels:
   app: rabbitmq
 template:
  metadata:
   labels:
    app: rabbitmq
  spec:
   serviceAccountName: rabbitmq-service-account
   initContainers:
    - name: config-loader
      image: docker.repos.divar.cloud/bhgedigital/envsubst
      imagePullPolicy: IfNotPresent
      command: [ '/bin/sh', '/tmp/config/init-config.sh' ]
      env:
       - name: RABBITMQ_ERLANG_COOKIE
         valueFrom:
          secretKeyRef:
           name: rabbitmq-secret
           key: RABBITMQ_ERLANG_COOKIE
       - name: RABBITMQ_ADMIN_PASSWORD
         valueFrom:
          secretKeyRef:
           name: rabbitmq-secret
           key: RABBITMQ_ADMIN_PASSWORD
      volumeMounts:
       - name: data
         mountPath: /var/lib/rabbitmq/
         readOnly: false
       - name: config
         mountPath: /tmp/config/
         readOnly: false
       - name: config-file
         mountPath: /etc/rabbitmq/
   containers:
    - name: rabbitmq
      image: docker.repos.divar.cloud/rabbitmq:3.9.14-management
      ports:
       - containerPort: 4369
         name: discovery
       - containerPort: 5672
         name: amqp
       - containerPort: 15672
         name: management
      env:
       - name: RABBIT_POD_NAME
         valueFrom:
          fieldRef:
           apiVersion: v1
           fieldPath: metadata.name
       - name: RABBIT_POD_NAMESPACE
         valueFrom:
          fieldRef:
           fieldPath: metadata.namespace
       - name: RABBITMQ_NODENAME
         value: rabbit@$(RABBIT_POD_NAME).rabbitmq-headless.$(RABBIT_POD_NAMESPACE).svc.cluster.local
       - name: RABBITMQ_USE_LONGNAME
         value: "true"
       - name: K8S_HOSTNAME_SUFFIX
         value: .rabbitmq-headless.$(RABBIT_POD_NAMESPACE).svc.cluster.local
      volumeMounts:
       - name: data
         mountPath: /var/lib/rabbitmq/
         readOnly: false
       - name: config-file
         mountPath: /etc/rabbitmq/
   volumes:
    - name: config-file
      emptyDir: { }
    - name: config
      configMap:
       name: rabbitmq-config
       defaultMode: 0755
 volumeClaimTemplates:
  - metadata:
     name: data
    spec:
     accessModes: [ "ReadWriteOnce" ]
     storageClassName: standard
     resources:
      requests:
       storage: 2Gi
```
<div dir="rtl">

از ابتدای فایل شروع می‌کنیم. اولش که یه سری اسم می‌دیم به پاد‌ها و بعدش پادمون به طور کلی دو تا کانتینر داره. اول یه کانتینر اولیه میاد که صرفا قراره اون اسکریپتی که توی کانفیگ مپ نوشتیم رو اجرا کنه. این کانتینر محتوای کانفیگ رو که توی فولدر `tmp/config/` قرار داره رو کپی می‌کنه به فولدر `etc/rabbitmq/` که این فولدر هم با یه مانت `emptyDir` به کانتینر اصلی‌مون وصل می‌شه. نکات بعدی این هست که ما در اینجا دوباره به سرویس `rabbitmq-headless` اشاره کردیم. این کار برای این هست که پاد‌هامون هاست‌نیم‌های استیبلی بگیرن که توی دی‌ان‌اس کوبر ریسالو بشن. در این مورد [اینجا](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#stable-network-id) بیشتر بخونید.

در ادامه سرویس اکانتی که ساختیم هم توی استیت فول ست نوشتیم که دسترسی‌ای که قبلا صحبت کردیم به پادامون بدیم.

پورت‌هایی که از کانتینر اصلی به بیرون دادیم همه پور‌ت‌های دیفالت هستند. پورت دیسکاوری برای ارتباط پاد‌ها به هم دیگه هست. پورت `ampq` برای ارتباط برنامه‌ها با کلاستر ربیت است و از طریق این پورت می‌تونن پیام‌ها رو بنویسینند یا بخوننند. پورت منیجمنت هم یه ui بهمون می‌ده.

در نهایت یک تمپلیت برای ساخت pvc در آخر فایل هست که برای این هست که هر کدوم از پاد‌ها استوریج مجزای خودشون رو داشته باشن و در هنگام اسکیل کرد هم به راحتی تعداد بیشتری ساخته بشه. حواستون باشه که `storageClassName` رو باید با توجه به کلاس‌هایی که توی کلاسترتون هست قرار بدید.

در ادامه می‌خواییم سرویس‌های مورد نیاز رو تعریف کنیم. اول باید دقت کنیم برای پورت دیسکاوری نیاز به سرویس هدلس داریم. در سرویس عادی‌ای که می‌سازیم ترافیک ورودی به سرویس بین پاد‌های پشت سرویس لودبالانس می‌شه. در سرویس هدلس ما دیگه یه هاست واحد نداریم و هر پاد یک هاست جدا می‌گیره. برای همین از سرویس هدلس استفاده می‌کنیم برای پورت دیسکاوری ولی دو پورت دیگه رو پشت یه سرویس عادی می‌ذاریم که ترافیک بین پاد‌هامون پخش بشه.

</div>

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    microservices: rabbitmq-headless
  name: rabbitmq-headless
spec:
  clusterIP: None
  ports:
    - port: 4369
      targetPort: 4369
      name: discovery
  selector:
    app: rabbitmq
---
apiVersion: v1
kind: Service
metadata:
  labels:
    microservices: rabbitmq
  name: rabbitmq
spec:
  ports:
    - port: 5672
      targetPort: 5672
      name: amqp
    - port: 15672
      targetPort: 15672
      name: management
  selector:
    app: rabbitmq
```

<div dir="rtl">

حال در ادامه مقداری در مورد نوع استفاده آن و یک کد سمپل در گولنگ می‌نویسیم. به طور کلی این موضوع برمی‌گردد به اجرای تسک‌های async. فرض کنید شما یک api وب دارید و این api شما در حین انجام کار خود یک کار سنگینی انجام می‌دهد. ممکن است این کار یک کار محاسباتی زیاد یا آپلود و دانلود یک فایل حجیم یا کال کردن یک سرویس دیگر باشد. در چنین مواقعی ممکن است به خاطر حجیم بود این کار api شما تایم‌اوت شود و شما بنابر لاجیک برنامه‌تان متوجه می‌شوید که لازم نیست که این تسک را در همان لحظه انجام دهید. می‌توانید ابتدا یک جواب به کاربر بدهید که ما داریم درخواست شما را انجام می‌دهیم و اون را به صفحه‌ای هدایت کنیم که وقتی جواب تسک حاضر شد در آنجا نتیجه را ببیند. در بعضی مواقع حتی دیدن نتیجه هم اهمیت چندانی ندارد. تصور کنید که یک سیستم احراز هویت داریم که از طریق پیامک زدن یا ایمیل زدن کار می‌کند. کاربر پس از وارد کردن نام کاربری منتظر پیامک خود است. ما می‌توانیم این پیامک را به صورت async در آیند‌ی نزدیک ارسال کنیم. راه حل اولیه شاید این باشد که این تسک‌ها در همون وب سرور پس از جواب دادن به کاربر با استفاده از ترد‌ها و ساختارهای async انجام دهیم.

در نگاه اول این روش می‌تواند کارا باشد اما در ادامه به مشکلاتی از آن برمی‌خوریم. مثلا فرض کنید این سرویس بر روی ریسورس منیجر کوبرنتیز قرار دارد و یکی از پاد‌های در حال اجرا متوقف  می‌شود. حال چه رخ می‌دهد؟ همه‌ی آن تسک‌های asyncdای که به کاربر قول داده بودیم انجام دهیم دیگر انجام نمی‌شود. یا حتی ممکن است با این روش لود پاد‌های ما زیاد شود دوست داشته باشیم برای این تسک‌های سنگین ریت‌لیمیتی تعیین کنیم و ریسورس جداگانه‌ای به آنها اختصاص دهیم. به طور کلی خواسته‌مون اینه که با افزایش این تسک‌های بزرگ و انباشه شدنشون روی هم دیگه سیستمی که داریم افکت نشه و بتونه همچنان سرویس خودش رو بده و این درخواست‌های رو در یک صف نگه داره که سر فرصت انجام بده و از طرفی دیگر این تضمین رو داشته باشیم که با کشته شدن یک پاد سیستم ما چیزی رو گم نمی‌کنه.

اینجا هست که ربیت‌ام‌کیو وارد ماجرا می‌شود. ربیت‌ام‌کیو به خاطر طراحی ساختاری‌ای که داره می‌تونه HA بودن رو تضمین کنه و با رپلیکیت کردن پیام‌ها تضمین می‌کنه که با کشته شدن یک پاد همچنان می‌تونه سرویس بده و خیلی راحت قابلیت اسکیل شدن هم داره. برای اسکیل کردنش یه ایده‌ی کلی هست که چند صف مشابه هم درست کنیم و داده‌هامون رو بین این صف‌ها پخش کنیم. از طرفی دیگر می‌تواد پاد‌های بیشتری به کلاستر ربیت اضافه کرد و این صف‌ها بین پاد‌های کلاستر ربیت پخش می‌شوند و عملا داده‌ی ما بین پاد‌های مختلف پخش می‌شود.

از طرفی دیگر ربیت‌ام‌کیو با ساختار اکنالجی که داره بعد از اینکه از مصرف پیام به صورت درستی مطمئن شد اون رو پاک می‌کنه. این خاصیت باعث می‌شه تضمین بشه هر پیامی دقیقا یک بار مصرف شده.

در انتها می‌خواییم یه کد سمپل بزنیم. برای این کار لایبرری celery بسیار محبوب است. این لایبرری در پایتون پیاده‌سازی نسبتا کاملی با فیچر‌های خوبی داره و نسخه‌ای از آن هم برای گولنگ پیاده شده که می‌خواهیم از آن استفاده کنیم.

ابتدا یک کد ساده برای ورکر آن در فایل 
`example/worker.go`
می‌نویسیم. در این کد نکته اول هاست مربوط به ربیت هست که فرمتی به شکل 
`amqp://user:password@host:port/vhost`
داره. بعد از اون ما باید اینترفیس `CeleryTask` رو برای تابع خودمون پیاده کنیم. در نهایت می‌تونیم در کدی که در فایل کلاینت قرار دارد یا تسک را فراخوانی کنیم.

برای اجرای ورکر کامند زیر را اجرا می‌کنیم.

```shell
go run main.go worker
```

و برای اجرای کلاینت تستی هم کامند زیر کمک می‌کنه.

```shell
go run main.go client
```


</div>

